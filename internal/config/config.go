package config

import "os"

type Config struct {
	BaseURL       string
	ServerAddress string
}

func Get() *Config {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	baseUrl := os.Getenv("BASE_URL")

	if serverAddress == "" {
		serverAddress = "0.0.0.0:8080"
	}
	if baseUrl == "" {
		baseUrl = "http://localhost:8080"
	}

	return &Config{
		BaseURL:       baseUrl,
		ServerAddress: serverAddress,
	}
}
