package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CookieSecure bool
	JWTSecret    string
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
		CookieSecure: os.Getenv("GO_ENV") != "dev",
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}
	fmt.Printf("AUTH: %v", AppConfig)
}
