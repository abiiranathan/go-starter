package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Application configuration struct populated from yaml or environment file.
/*
Corresponding Environment variables:

DATABASE_URL="postgres://user:password@localhost:5432/dbname"

PORT=8080

SECRET_KEY=some super secure key

DEBUG=true

CORS_ALLOWED_ORIGINS=

REDIS_CLIENT_ADDR=
*/
type Config struct {
	DatabaseURL        string   `yaml:"database_url"`
	Port               string   `yaml:"port"`
	SecretKey          string   `yaml:"secret_key"`
	Debug              bool     `yaml:"debug"`
	CorsAllowedOrigins []string `yaml:"cors_allowed_origins"`
	RedisClientAddr    string   `yaml:"redis_client_addr"`
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv(envFiles ...string) (*Config, error) {
	err := godotenv.Load(envFiles...)
	if err != nil {
		return nil, err
	}

	var debugEnv bool = true
	debug, found := os.LookupEnv("DEBUG")
	if found {
		debugEnv = (debug == "1" || debug == "true")
	}

	c := &Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		Port:               os.Getenv("PORT"),
		SecretKey:          os.Getenv("SECRET_KEY"),
		Debug:              debugEnv,
		CorsAllowedOrigins: strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		RedisClientAddr:    os.Getenv("REDIS_CLIENT_ADDR"),
	}

	setDefaults(c)
	return c, nil
}

// LoadFromYAML loads configuration from a YAML file.
func LoadFromYAML(yamlFile string) (*Config, error) {
	c := &Config{}
	err := loadYAML(yamlFile, c)
	if err != nil {
		return nil, err
	}
	setDefaults(c)
	return c, nil
}

func loadYAML(yamlFile string, c *Config) error {
	f, err := os.Open(yamlFile)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}

	return nil
}

func setDefaults(c *Config) {
	if c.Port == "" {
		c.Port = "8080"
	}

	if c.RedisClientAddr == "" {
		c.RedisClientAddr = "127.0.0.1:6379"
	}
}

// Load the yaml/.env configuration based on command-line flags.
func Load() (*Config, error) {
	var configPath string
	var env string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&env, "env", "", "path to env file")
	flag.Parse()

	if configPath != "" && env != "" {
		return nil, fmt.Errorf("only one of -config or -env flags can be provided")
	}

	if configPath != "" {
		c, err := LoadFromYAML(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}
		return c, nil
	}

	if env != "" {
		c, err := LoadFromEnv(env)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}
		return c, nil
	}

	// If .env is available, load from it
	stat, err := os.Stat(".env")
	if err == nil && !stat.IsDir() {
		c, err := LoadFromEnv(".env")
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}
		return c, nil
	}
	return nil, fmt.Errorf("no configuration provided")
}
