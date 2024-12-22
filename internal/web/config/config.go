package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoEnv string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	} else {
		log.Println(".env file initialized.")
	}

	AppConfig = &Config{
		GoEnv: os.Getenv("GO_ENV"),
	}
}
