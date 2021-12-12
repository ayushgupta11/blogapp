package user

import (
	"blogapp/helper"
	"blogapp/helper/hash"
	"blogapp/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Connection mongoDB with helper class
var collection = helper.ConnectToUsers()

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.Password = hash.HashAndSalt([]byte(user.Password))
	layout := "2006-01-02"
	str := user.Dob
	t, err := time.Parse(layout, str)
	fmt.Println(t)
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte("Bad request. Please enter date in YYYY-MM-DD format."))
	} else {
		result, err := collection.InsertOne(context.TODO(), user)

		if err != nil {
			helper.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var user models.User
	var newUser models.User
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&user)
	if user.Name != "" && user.Email != "" && user.Password != "" {
		update := bson.D{
			{"$set", bson.D{
				{"name", user.Name},
				{"email", user.Email},
				{"dob", user.Dob},
				{"phone", user.Phone},
				{"password", hash.HashAndSalt([]byte(user.Password))},
			}},
		}
		err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&newUser)
		if err != nil {
			helper.GetError(err, w)
			return
		}
		newUser.ID = id
		newUser.Password = ""
		json.NewEncoder(w).Encode(newUser)
	} else {
		w.WriteHeader(422)
		w.Write([]byte("Bad request"))
		return
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inputUser models.User
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&inputUser)
	fmt.Println(inputUser)
	filter := bson.M{"username": inputUser.Username}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return
	} else {
		var pwdMatch = hash.ComparePasswords(user.Password, []byte(inputUser.Password))
		user.Password = ""
		if !pwdMatch {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	user.Password = ""
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return
	}

	json.NewEncoder(w).Encode(user)
}
