package main

import (
	"log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

var (
	Port	string
)

// handlers
func loginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
}

// config, etc.
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed while loading config")
	}
	Port = os.Getenv("PORT")
}

func main() {
	r := mux.NewRouter()
	// routes
	r.HandleFunc("/user/login", loginUser).Methods("POST")
	r.HandleFunc("/user/register", registerUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+Port, r))
}