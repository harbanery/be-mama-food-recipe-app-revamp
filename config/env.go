package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvLoad struct {
	APP_NAME      string

	WEB_PREFORK   string
	WEB_PORT       string

	CORS_ALLOW_ORIGINS string
	CORS_ALLOW_METHODS string
	CORS_ALLOW_HEADERS string
	CORS_EXPOSE_HEADERS string

	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASS string
	DB_NAME string
	DB_SSLMODE string
	DB_POOL_IDLE string
	DB_POOL_MAX string
	DB_POOL_LIFETIME string
}

func CheckEnv() *EnvLoad {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return &EnvLoad{
		APP_NAME:      os.Getenv("APP_NAME"),

		WEB_PREFORK:   os.Getenv("WEB_PREFORK"),
		WEB_PORT:      os.Getenv("WEB_PORT"),

		CORS_ALLOW_ORIGINS: os.Getenv("CORS_ALLOW_ORIGINS"),
		CORS_ALLOW_METHODS: os.Getenv("CORS_ALLOW_METHODS"),
		CORS_ALLOW_HEADERS: os.Getenv("CORS_ALLOW_HEADERS"),
		CORS_EXPOSE_HEADERS: os.Getenv("CORS_EXPOSE_HEADERS"),

		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
		DB_NAME: os.Getenv("DB_NAME"),
		DB_SSLMODE: os.Getenv("DB_SSLMODE"),
		DB_POOL_IDLE: os.Getenv("DB_POOL_IDLE"),
		DB_POOL_MAX: os.Getenv("DB_POOL_MAX"),
		DB_POOL_LIFETIME: os.Getenv("DB_POOL_LIFETIME"),
	}
}

func GetPort() string {
	var env EnvLoad

	port := env.WEB_PORT
	if port == "" {
		port = "3003"
	}

	return "0.0.0.0:" + port
}