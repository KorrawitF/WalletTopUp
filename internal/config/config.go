package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      Server
	Database    Database
	CacheClient CacheClient
}

type Server struct {
	Host string
	Port string
	Mode string
}

type Database struct {
	Host string
	Port string
	User string
	Pass string
	Name string
	Mode string
}

type CacheClient struct {
	Enabled bool
	Host    string
	Port    string
	Pass    string
	DB      int
}

func Load() *Config {
	godotenv.Load(".env")
	return &Config{
		Server: Server{
			Host: lookUpEnv[string]("SERVER_HOST", "0.0.0.0"),
			Port: lookUpEnv[string]("SERVER_PORT", "8000"),
			Mode: lookUpEnv[string]("SERVER_MODE", "release"),
		},
		Database: Database{
			Host: lookUpEnv[string]("DB_HOST", "localhost"),
			Port: lookUpEnv[string]("DB_PORT", "5432"),
			User: lookUpEnv[string]("DB_USER", "root"),
			Pass: lookUpEnv[string]("DB_PASS", "pass"),
			Name: lookUpEnv[string]("DB_NAME", "postgres"),
			Mode: lookUpEnv[string]("DB_MODE", "disable"),
		},
		CacheClient: CacheClient{
			Enabled: lookUpEnv[bool]("CACHE_ENABLED", "false"),
			Host:    lookUpEnv[string]("CACHE_HOST", "localhost"),
			Port:    lookUpEnv[string]("CACHE_PORT", "6379"),
			Pass:    lookUpEnv[string]("CACHE_PASS", ""),
			DB:      lookUpEnv[int]("CACHE_DB", "0"),
		},
	}
}

func lookUpEnv[T string | int | bool](key string, fallback ...string) T {
	val, exists := os.LookupEnv(key)
	if !exists {
		if len(fallback) == 0 {
			log.Fatalf("Env key: %s is required but not set.", key)
		}
		val = fallback[0]
	}

	var result T
	switch parsed := any(&result).(type) {
	case *bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatal(err)
		}
		*parsed = boolVal
	case *string:
		*parsed = val
	case *int:
		intVal, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Env key: %s must be an integer, got %q.", key, val)
		}
		*parsed = intVal
	}
	return result
}
