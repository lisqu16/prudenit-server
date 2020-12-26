package config

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	Port		string
	Address		string
	Name		string
	User		string
	Password	string
)

func Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// variables
	Port = os.Getenv("PORT")
	Address = os.Getenv("DBADDRESS")
	Name = os.Getenv("DBNAME")
	User = os.Getenv("DBUSER")
	Password = os.Getenv("DBPASS")
	return nil
}