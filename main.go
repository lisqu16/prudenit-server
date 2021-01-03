package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/lisqu16/prudenit-server/user"
	"github.com/lisqu16/prudenit-server/config"
	db "github.com/lisqu16/prudenit-server/database"
)

// config, etc.
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed while loading config:\n%s", err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatalf("Failed while connecting with database:\n%s", err)
	}
}

func main() {
	r := mux.NewRouter()
	// routes
	r.HandleFunc("/user/login", user.Login).Methods("POST")
	r.HandleFunc("/user/register", user.Register).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}