package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	KeycloakURL string // Например: http://192.168.64.2:8080
	ClientID    string // notes-api
	Realm       string // my-project
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Предупреждение: .env файл не найден")
	}

	return &Config{
		KeycloakURL: os.Getenv("KEYCLOAK_URL"),
		ClientID:    os.Getenv("KEYCLOAK_CLIENT_ID"),
		Realm:       os.Getenv("KEYCLOAK_REALM"),
	}
}
