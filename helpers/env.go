package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Key struct {
	UploadThing_Key string
	DbUser          string
	DbPassword      string
	DbName          string
	DbHost          string
	DbPort          string
	SslMode         string
}

func init() {
	err := godotenv.Load(".env")
	if err!= nil {
		fmt.Println("missing variables continue")
	}
}

func FetchEnv() Key {

	var key Key

	key.DbUser = os.Getenv("DB_USER")
	key.DbPassword = os.Getenv("DB_PASSWORD")
	key.DbName = os.Getenv("DB_NAME")
	key.DbHost = os.Getenv("DB_HOST")
	key.DbPort = os.Getenv("DB_PORT")
	key.SslMode = os.Getenv("SSL_MODE")
	key.UploadThing_Key = os.Getenv("UPLOAD_THING_KEY")

	if key.SslMode == "" {
		key.SslMode = "disable"
	}

	// Validate required environment variables
	if key.DbUser == "" || key.DbPassword == "" || key.DbName == "" || key.DbHost == "" {
		log.Fatal("Missing required environment variables")
	}

	return key
}
