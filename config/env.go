package config

import "os"

type Env struct {
	ClientID string
}

func LoadEnv() *Env {
	return &Env{
		ClientID: getEnv("CLIENT_ID", "Ov23liM0BAkzFlF1II7n"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
