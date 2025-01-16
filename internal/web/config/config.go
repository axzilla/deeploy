package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	CookieSecure bool
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
		// CookieSecure: os.Getenv("GO_ENV") != "dev",
		CookieSecure: false, // INFO: false for now because to many http/https cookie related issues (e.g login) while testing
	}
}
