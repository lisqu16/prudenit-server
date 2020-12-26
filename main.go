package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
)

func loginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("login")
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	log.Println("register")
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	r := mux.NewRouter()
	// routes
	r.HandleFunc("/user/login", loginUser).Methods("POST")
	r.HandleFunc("/user/register", registerUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}