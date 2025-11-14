package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

var JwtKey []byte

func LoadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    JwtKey = []byte(os.Getenv("JWT_SECRET"))
}
