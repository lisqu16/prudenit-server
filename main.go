package main

import (
	"log"
	"fmt"
	"regexp"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"github.com/lisqu16/prudenit-server/config"
	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var (
	s *rethink.Session
)

type User struct {
	Email		string `json:"email"`
	Username	string `json:"username"`
}

type Response struct {
	Ok		bool `json:"ok"`
	Msgs 	[]string `json:"messages"`	
	Data 	map[string]interface{} `json:"data"`
}

// handlers
func loginUser(w http.ResponseWriter, r *http.Request) {
	var ok bool = true
	var msgs []string

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	cursor, _ := rethink.Table("users").Filter(rethink.Row.Field("email").Eq(email)).Run(s)
	var res map[string]string
	err := cursor.One(&res)
	if err != rethink.ErrEmptyResult {
		hashedPassword := res["password"]
		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			ok = false
			msgs = append(msgs, "wrongEmailOrPassword")
		}
	} else {
		ok = false
		msgs = append(msgs, "wrongEmailOrPassword")
	}
	defer cursor.Close()

	w.Header().Set("Content-Type", "application/json; charset=utf8")

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{ok, msgs, make(map[string]interface{})})
		return
	}

	w.WriteHeader(http.StatusOK)

	var data = make(map[string]interface{})
	data["user"] = &User{res["email"], res["username"]}
	json.NewEncoder(w).Encode(Response{ok, msgs, data})
	return
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var ok bool = true
	var msgs []string

	r.ParseForm()
	email := r.Form.Get("email")
	username := r.Form.Get("username")
	pass := r.Form.Get("password")
	repass := r.Form.Get("repassword")

	// regex stolen from another website hehe
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	
	match := emailRegex.MatchString(email)
	if !match {
		msgs = append(msgs, "wrongEmail")
	}

	if len(username) < 2 || len(username) > 24 {
		msgs = append(msgs, "wrongUsernameLength")
	}

	if pass != repass {
		msgs = append(msgs, "differentPasswords")
	}
	if len(pass) <= 7 || len(pass) > 64 {
		msgs = append(msgs, "tooShortOrLongPassword")
	}

	if match {
		cursor, _ := rethink.Table("users").Filter(rethink.Row.Field("email").Eq(email)).Run(s)
		var res interface{}
		err := cursor.One(res)
		if err != rethink.ErrEmptyResult {
			msgs = append(msgs, "alreadyAssignedEmail")
		}
		defer cursor.Close()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	if len(msgs) > 0 {
		ok = false
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(Response{ok, msgs, make(map[string]interface{})})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	var data = make(map[string]string)
	data["email"] = email
	data["username"] = username
	data["password"] = string(hashedPassword)
	rethink.Table("users").Insert(data).Run(s)
	
	w.WriteHeader(http.StatusCreated)

	var rdata = make(map[string]interface{})
	rdata["user"] = &User{data["email"], data["username"]}
	json.NewEncoder(w).Encode(Response{ok, msgs, rdata})
	return
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