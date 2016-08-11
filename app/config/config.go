package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port       string
	MongoDBUrl string
	GinEnv     string
}

func Load(envFile string) (*Config, error) {
	if envFile == "" {
		envFile = ".env"
	}

	fmt.Printf("Loading %s file\n", envFile)

	godotenv.Load(envFile)

	mongoDBUrl := os.Getenv("MONGODB_URL")
	if mongoDBUrl == "" {
		return nil, fmt.Errorf("Missing MONGODB_URL")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	ginEnv := os.Getenv("GIN_ENV")
	if ginEnv == "" {
		ginEnv = "development"
	}

	return &Config{
		port,
		mongoDBUrl,
		ginEnv,
	}, nil
}
