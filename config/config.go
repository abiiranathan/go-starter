package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabaseURL        string   `yaml:"database_url"`
	Port               string   `yaml:"port"`
	SecretKey          string   `yaml:"secret_key"`
	Debug              bool     `yaml:"debug"`
	CorsAllowedOrigins []string `yaml:"cors_allowed_origins"`
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv(envFiles ...string) (*Config, error) {
	err := godotenv.Load(envFiles...)
	if err != nil {
		return nil, err
	}

	c := &Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		Port:               os.Getenv("PORT"),
		SecretKey:          os.Getenv("SECRET_KEY"),
		Debug:              os.Getenv("DEBUG") == "true",
		CorsAllowedOrigins: strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
	}

	return c, nil
}

// LoadFromYAML loads configuration from a YAML file.
func LoadFromYAML(yamlFile string) (*Config, error) {
	c := &Config{}
	err := loadYAML(yamlFile, c)
	if err != nil {
		return nil, err
	}

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
