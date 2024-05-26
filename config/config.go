package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	defaultTimeout = 10
)

// Config is a container for all the needed app configuration.
type Config struct {
	Server   Server
	Database Database
}

// Server holds the server configuration.
type Server struct {
	Host    string        `default:"localhost" env:"SERVER_HOST"`
	Port    string        `default:"8080" env:"SERVER_PORT"`
	Timeout time.Duration `default:"10" env:"SERVER_TIMEOUT"`
}

// Database holds the database configuration.
type Database struct {
	Type string `default:"in-memory" env:"DATABASE_TYPE"`
}

// New initialize the config.
func New() *Config {
	return &Config{
		Server: Server{
			Host:    getOrDefault("SERVER_HOST", "localhost"),
			Port:    getOrDefault("SERVER_PORT", "8080"),
			Timeout: defaultTimeout,
		},
		Database: Database{
			Type: getOrDefault("DATABASE_TYPE", "in-memory"),
		},
	}
}

// getOrDefault gets the value from environment if not returns the default value.
func getOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}

// InitLogger initializes the logging.
func InitLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:       true,
		DisableHTMLEscape: true,
		TimestampFormat:   time.RFC3339,
	})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
