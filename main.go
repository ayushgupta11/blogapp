package main

import (
	"log"
	"net/http"

	"blogapp/controllers/posts"
	"blogapp/controllers/user"
	"blogapp/helper"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/user", user.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/login", user.LoginUser).Methods("POST")
	r.HandleFunc("/api/user/{id}", user.GetUser).Methods("GET")
	r.HandleFunc("/api/user/{id}", user.EditUser).Methods("PUT")
	r.HandleFunc("/api/user/search", user.SearchUser).Methods("POST")

	r.HandleFunc("/api/posts", posts.GetPosts).Methods("GET")
	r.HandleFunc("/api/posts", posts.CreatePost).Methods("POST")
	r.HandleFunc("/api/posts/{id}", posts.DeletePost).Methods("DELETE")
	r.HandleFunc("/api/posts/{id}", posts.EditPost).Methods("PUT")

	config := helper.GetConfiguration()
	log.Fatal(http.ListenAndServe(config.Port, r))

}
