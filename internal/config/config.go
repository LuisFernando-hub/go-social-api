package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)


type Config struct {
	Host string
	DataBase string
	User string
	Password string
	Port string
	JWTSecret string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found, using enviroment variables")
	}

	var config *Config = &Config{
		Host: os.Getenv("DB_HOST"),
		User: os.Getenv("DB_USER"),
		DataBase: os.Getenv("DB_DATABASE"),
		Password: os.Getenv("DB_PASSWORD"),
		Port: os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return config, nil
}