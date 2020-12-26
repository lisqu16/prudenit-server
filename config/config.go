package config

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	Port	string
)

func Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	Port = os.Getenv("PORT")
	return nil
}