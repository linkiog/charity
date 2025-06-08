package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	JWTSecret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("end not founf or could be loaded")
	}
	return &Config{
		DBUrl:     os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
