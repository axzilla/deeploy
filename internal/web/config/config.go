package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CookieSecure  bool
	AssetCaching  bool
	UseEmbeddedFS bool
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
		CookieSecure:  os.Getenv("GO_ENV") != "dev",
		AssetCaching:  os.Getenv("GO_ENV") != "dev",
		UseEmbeddedFS: os.Getenv("GO_ENV") != "dev",
	}
}
