package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/lisqu16/prudenit-server-go/config"
	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var (
	s *rethink.Session
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

	s, err = rethink.Connect(rethink.ConnectOpts{
		Address: config.Address,
		Database: config.Name,
		Username: config.User,
		Password: config.Password,
	})
	if err != nil {
		log.Fatal(fmt.Printf("Failed while connecting with database\n%s", err))
	}
}

func main() {
	r := mux.NewRouter()
	// routes
	r.HandleFunc("/user/login", loginUser).Methods("POST")
	r.HandleFunc("/user/register", registerUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}