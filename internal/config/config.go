package config

import (
	"net"
	"net/url"
	"os"
)

// Config contains runtime configuration pulled from environment variables.
type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
}

// Load reads configuration from environment variables, falling back to safe defaults.
func Load() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("POSTGRES_USER", "postgres"),
		DBPass:    getEnv("POSTGRES_PASSWORD", ""),
		DBName:    getEnv("POSTGRES_DB", "todos"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
	}
}

// Addr returns the HTTP listen address.
func (c Config) Addr() string {
	return ":" + c.Port
}

// DatabaseURL builds a PostgreSQL connection string.
func (c Config) DatabaseURL() string {
	host := net.JoinHostPort(c.DBHost, c.DBPort)
	u := url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   "/" + c.DBName,
	}

	if c.DBPass != "" {
		u.User = url.UserPassword(c.DBUser, c.DBPass)
	} else {
		u.User = url.User(c.DBUser)
	}

	q := u.Query()
	q.Set("sslmode", c.DBSSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
