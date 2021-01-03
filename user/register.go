package user

import (
	"regexp"
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	db "github.com/lisqu16/prudenit-server/database"
	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var ok bool = true
	var msgs []string

	r.ParseForm()
	email := r.Form.Get("email")
	username := r.Form.Get("username")
	pass := r.Form.Get("password")
	repass := r.Form.Get("repassword")
	
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
		cursor, _ := rethink.Table("users").Filter(rethink.Row.Field("email").Eq(email)).Run(db.Session)
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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	var data = make(map[string]string)
	data["email"] = email
	data["username"] = username
	data["password"] = string(hashedPassword)
	rethink.Table("users").Insert(data).Run(db.Session)
	
	w.WriteHeader(http.StatusCreated)

	var rdata = make(map[string]interface{})
	rdata["user"] = &User{data["email"], data["username"]}
	json.NewEncoder(w).Encode(Response{ok, msgs, rdata})
	return
}