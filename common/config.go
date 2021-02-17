package common

import "os"

//Config is general struct for configurations
type Config struct {
	DBConfig  DBConfig
	AppConfig AppConfig
}

//DBConfig contains database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

//AppConfig contains application configuration
type AppConfig struct {
	Port string
}

// FromEnv returns Config structure filled from environment variables
func FromEnv() Config {
	dbc := DBConfig{
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	}
	ac := AppConfig{
		os.Getenv("APP_PORT"),
	}
	return Config{dbc, ac}
}
