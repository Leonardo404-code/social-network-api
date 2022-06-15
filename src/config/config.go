package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	StringConnection = ""

	Port = 0

	SecretKey []byte
)

// Load configure all global vars and environments variables
func Load() {
	var error error

	if error = godotenv.Load(); error != nil {
		log.Fatalf("error in load environment variable: %v", error)
	}

	Port, error = strconv.Atoi(os.Getenv("API_PORT"))

	if error != nil {
		Port = 9000
	}

	user := os.Getenv("DB_USUARIO")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NOME")
	secret := os.Getenv("SECRET")

	StringConnection = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", user, password, dbName)

	SecretKey = []byte(secret)
}
