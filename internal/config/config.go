// Package config provides application configuration management.
// It handles loading and parsing configuration files (YAML) with environment variable substitution,
// and supports loading environment variables from .env files.
// The package supports command-line flags for specifying config and .env file locations.
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration structure.
// It contains all configurable parameters grouped by logical components.
type Config struct {
	HTTPServer HTTPServer `yaml:"http-server"`
	PgDSN      string     `yaml:"pg-dsn"`
}

// HTTPServer contains configuration parameters for the HTTP server.
type HTTPServer struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

// LoadConfig loads and parses the application configuration.
// It performs the following steps:
//  1. Parses command-line flags for config and .env file locations
//  2. Loads environment variables from the specified .env file (if exists)
//  3. Reads the configuration file (defaults to ./config/config.yaml)
//  4. Expands environment variables in the config file
//  5. Unmarshals the YAML content into the Config struct
//
// Returns:
//   - *Config: Loaded configuration object
//   - error: Any error that occurred during loading or parsing
//
// Note: The function uses os.Expand to substitute environment variables in the config file,
// allowing for dynamic configuration values.
func LoadConfig() (*Config, error) {
	envPathFlag := flag.String("env", ".env", "Path to .env file")
	configPathFlag := flag.String("config", "", "Path to config file")
	flag.Parse()

	if err := godotenv.Load(*envPathFlag); err != nil {
		fmt.Printf(".env file not found: %s\n", *envPathFlag)
	}

	configPath := *configPathFlag
	if configPath == "" {
		configPath = filepath.Join(".", "config", "config.yaml")
	}

	data, err := os.ReadFile(configPath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	processedData := os.Expand(string(data), os.Getenv)

	var cfg Config

	if err = yaml.Unmarshal([]byte(processedData), &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}
