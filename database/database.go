package database

import (
	"log"
	"github.com/lisqu16/prudenit-server/config"
	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var (
	Session *rethink.Session
)

func Connect() error {
	// session
	var err error
	Session, err = rethink.Connect(rethink.ConnectOpts{
		Address: config.Address,
		Username: config.User,
		Password: config.Password,
	})

	if err != nil {
		return err
	}

	// database
	var exists bool
	cursor, _ := rethink.DBList().Contains(config.Name).Run(Session)
	cursor.One(&exists)
	if !exists {
		rethink.DBCreate(config.Name).Run(Session)
		log.Printf("Created database \"%s\"", config.Name)
	}
	Session.Use(config.Name)

	// table
	var tableExist bool
	cursor, _= rethink.TableList().Contains("users").Run(Session)
	cursor.One(&tableExist)
	if !tableExist {
		rethink.TableCreate("users").Run(Session)
	}

	return nil
}