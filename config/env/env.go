package env

import (
    "github.com/joho/godotenv"
    "log"
    "os"
    "strconv"
)

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment")
    }
}

func GetEnv(key, fallback string) string {
    value := os.Getenv(key)
    if value == "" {
        return fallback
    }
    return value
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
