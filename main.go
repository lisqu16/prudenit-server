package main

import (
	"log"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/lisqu16/prudenit-server-go/config"
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
	err := config.Load()
	if err != nil {
		log.Fatal(fmt.Printf("Failed while loading config\n%s", err))
	}
}

func main() {
	r := mux.NewRouter()
	// routes
	r.HandleFunc("/user/login", loginUser).Methods("POST")
	r.HandleFunc("/user/register", registerUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}