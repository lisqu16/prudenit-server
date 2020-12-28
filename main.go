package main

import (
	"log"
	"fmt"
	"regexp"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"github.com/lisqu16/prudenit-server-go/config"
	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var (
	s *rethink.Session
)

type User struct {
	Email		string `json:"email"`
	Username	string `json:"username"`
}

// handlers
func loginUser(w http.ResponseWriter, r *http.Request) {
	var ok bool = true
	var msg string

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
			msg = "wrongEmailOrPassword"
		}
	} else {
		ok = false
		msg = "wrongEmailOrPassword"
	}
	defer cursor.Close()

	w.Header().Set("Content-Type", "application/json; charset=utf8")

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		type Fail struct {
			Ok		bool `json:"ok"`
			Msg		string `json:"message"`
		}
		json.NewEncoder(w).Encode(Fail{ok, msg})
		return
	}

	w.WriteHeader(http.StatusOK)
	type Success struct {
		Ok		bool `json:"ok"`
		User	*User `json:"user"`
	}
	json.NewEncoder(w).Encode(Success{ok, &User{res["email"], res["username"]}})
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
	
	if match := emailRegex.MatchString(email); !match {
		msgs = append(msgs, "wrongEmail")
	}

	if len(username) < 2 || len(username) > 24 {
		msgs = append(msgs, "wrongUsernameLength")
	}

	if pass != repass {
		msgs = append(msgs, "differentPasswords")
	}
	if len(pass) <= 7 {
		msgs = append(msgs, "tooShortPassword")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	if len(msgs) > 0 {
		ok = false
		type Fail struct {
			Ok		bool `json:"ok"`
			Msgs	[]string `json:"messages"`
		}
		json.NewEncoder(w).Encode(Fail{ok, msgs})
		return
	}
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