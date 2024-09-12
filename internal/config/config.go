package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host             string `yaml:"host" envconfig:"HOST"`
	Port             int    `yaml:"port" envconfig:"PORT"`
	DBPassword       string `yaml:"db_password" envconfig:"DB_PASSWORD"`
	CrtPath          string `yaml:"crt_path" envconfig:"CRT_PATH"`
	KeyPath          string `yaml:"key_path" envconfig:"KEY_PATH"`
	SecretKey        []byte `yaml:"secret_key" envconfig:"SECRET_KEY"`
	SecretRefreshKey []byte `yaml:"secret_refresh_key" envconfig:"SECRET_REFRESH_KEY"`
	CurrencyURL      string `yaml:"currency_url" envconfig:"CURRENCY_URL"`
}

func New() (*Config, error) {
	var cfg Config

	if err := loadConfig("config/config.yml", &cfg); err != nil {
		return nil, err
	}
	if err := processEnvironment(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c Config) GetDBURL() string {
	return fmt.Sprintf("postgres://cadvadmin:%s@%s:5432/cadvdb?sslmode=disable", c.DBPassword, c.Host)
}

func loadConfig(filename string, cfg *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", filename, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode YAML from config file %s: %w", filename, err)
	}

	return nil
}

func processEnvironment(cfg *Config) error {
	if err := envconfig.Process("", cfg); err != nil {
		return fmt.Errorf("failed to process environment variables: %w", err)
	}
	return nil
}
