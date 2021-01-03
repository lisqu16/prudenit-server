package user

import (
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	db "github.com/lisqu16/prudenit-server/database"
	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var ok bool = true
	var msgs []string

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	cursor, _ := rethink.Table("users").Filter(rethink.Row.Field("email").Eq(email)).Run(db.Session)
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
