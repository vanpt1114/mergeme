package config

import (
	"fmt"
	"os"
)

type Config struct {
	Gitlab			Gitlab
	Redis			Redis
	SlackToken		string
}

type Gitlab struct {
	URL		string
	Token	string
}

type Redis struct {
	Host		string
	Password	string
	DB			string
}

func mustHaveEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Sprintf("Environment variable %s does not exist", key))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Load() (c *Config) {
	c = &Config{
		Gitlab:     Gitlab{
			URL:   mustHaveEnv("GITLAB_URL"),
			Token: mustHaveEnv("GITLAB_TOKEN"),
		},
		Redis:      Redis{
			Host:     mustHaveEnv("REDIS_HOST"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnv("REDIS_DB", "0"),
		},
		SlackToken: mustHaveEnv("SLACK_TOKEN"),
	}
	return c
}
