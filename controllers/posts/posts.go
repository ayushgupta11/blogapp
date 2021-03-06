package posts

import (
	"blogapp/helper"
	"blogapp/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Connection mongoDB with helper class
var collection = helper.ConnectToPosts()

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var post models.Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	post.PostedOn = time.Now()
	var thumbnailUrl = "http://google.com/"
	if post.Thumbnail != "" {
		thumbnailUrl = post.Thumbnail
	}
	u, err := url.ParseRequestURI(thumbnailUrl)
	fmt.Println(u)
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte("Bad request"))
		return
	} else {
		result, err := collection.InsertOne(context.TODO(), post)
		if err != nil {
			w.WriteHeader(422)
			w.Write([]byte("Bad request"))
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var post models.Post
	var newPost models.Post
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&post)
	if post.Title != "" && post.Body != "" {
		update := bson.D{
			{"$set", bson.D{
				{"title", post.Title},
				{"body", post.Body},
				{"thumbnail", post.Thumbnail},
			}},
		}
		var thumbnailUrl = "http://google.com/"
		if post.Thumbnail != "" {
			thumbnailUrl = post.Thumbnail
		}
		u, err := url.ParseRequestURI(thumbnailUrl)
		fmt.Println(u)
		if err != nil {
			w.WriteHeader(422)
			w.Write([]byte("Bad request"))
			return
		} else {
			err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&newPost)
			if err != nil {
				helper.GetError(err, w)
				return
			}
			newPost.ID = id
			json.NewEncoder(w).Encode(newPost)
		}
	} else {
		w.WriteHeader(422)
		w.Write([]byte("Bad request"))
		return
	}
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []models.Post
	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{"postedOn", -1}})
	cur, err := collection.Find(context.TODO(), bson.M{}, findOptions)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var post models.Post
		err := cur.Decode(&post)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(posts)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
}
