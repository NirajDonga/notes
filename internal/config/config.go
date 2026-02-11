package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI    string
	MongoDbName string
	ServerPORT  string
}

func Load() (Config, error) {

	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("Failed to load .env")
	}

	mongoURI, err := exctractEnv("MONGODB_URI")
	if err != nil {
		return Config{}, err
	}

	mongoDB, err := exctractEnv("MONGODB_NAME")
	if err != nil {
		return Config{}, err
	}

	port, err := exctractEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	return Config{
		MongoURI:    mongoURI,
		MongoDbName: mongoDB,
		ServerPORT:  port,
	}, nil

}

func exctractEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing %v", key)
	}

	return val, nil
}
