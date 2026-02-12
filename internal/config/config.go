package config

import (
	"os"
)

type Config struct {
    AppPort        string
    DBHost         string
    DBPort         string
    DBUser         string
    DBPassword     string
    DBName         string
    RedisHost      string
    RedisPort      string
    RabbitMQURL    string
}

func LoadConfig() *Config {
    return &Config{
        AppPort:        getEnv("APP_PORT", ":8080"),
        DBHost:         getEnv("DB_HOST", "localhost"),
        DBPort:         getEnv("DB_PORT", "5432"),
        DBUser:         getEnv("DB_USER", "user"),
        DBPassword:     getEnv("DB_PASSWORD", "password"),
        DBName:         getEnv("DB_NAME", "property_db"),
        RedisHost:      getEnv("REDIS_HOST", "localhost"),
        RedisPort:      getEnv("REDIS_PORT", "6379"),
        RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://user:password@localhost:5672/"),
    }
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
