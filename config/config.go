package config

import "os"

// DBConfig  Database configuration structure
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// FromEnv returns DBConfig structure filled from environment variables
func FromEnv() DBConfig {
	dbc := DBConfig{
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	}
	return dbc
}
