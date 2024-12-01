package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/wachrusz/Back-End-API/pkg/rabbit"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

type Config struct {
	Host                string        `yaml:"host"`
	Port                int           `yaml:"port"`
	DBPassword          string        `yaml:"db_password"`
	CrtPath             string        `yaml:"crt_path"`
	KeyPath             string        `yaml:"key_path"`
	SecretKey           []byte        `yaml:"secret_key"`
	SecretRefreshKey    []byte        `yaml:"secret_refresh_key"`
	CurrencyURL         string        `yaml:"currency_url"`
	Rabbit              rabbit.Config `yaml:"rabbit"`
	AccessTokenLifetime int           `yaml:"access_token_dur_minutes"`
	RateLimitPerSecond  int           `yaml:"rate_limit_per_second"`
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
	err := godotenv.Load("secret/.env")
	if err != nil {
		return err
	}
	if host, exists := os.LookupEnv("HOST"); exists {
		cfg.Host = host
	}
	if portStr, exists := os.LookupEnv("PORT"); exists {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid PORT value: %w", err)
		}
		cfg.Port = port
	}
	if dbPassword, exists := os.LookupEnv("DB_PASSWORD"); exists {
		cfg.DBPassword = dbPassword
	}
	if secretKey, exists := os.LookupEnv("SECRET_KEY"); exists {
		cfg.SecretKey = []byte(secretKey)
	}
	if secretRefreshKey, exists := os.LookupEnv("SECRET_REFRESH_KEY"); exists {
		cfg.SecretRefreshKey = []byte(secretRefreshKey)
	}
	if currencyURL, exists := os.LookupEnv("CURRENCY_URL"); exists {
		cfg.CurrencyURL = currencyURL
	}

	if rabbitUrl, exists := os.LookupEnv("RABBIT_URL"); exists {
		cfg.Rabbit.URL = rabbitUrl
	}

	if tokenLifetimeStr, exists := os.LookupEnv("ACCESS_LIFETIME"); exists {
		tokenLifetime, err := strconv.Atoi(tokenLifetimeStr)
		if err != nil {
			return fmt.Errorf("invalid access token lifetime value: %w", err)
		}
		cfg.AccessTokenLifetime = tokenLifetime
	}

	if rateStr, exists := os.LookupEnv("RATE_LIMIT"); exists {
		rate, err := strconv.Atoi(rateStr)
		if err != nil {
			return fmt.Errorf("invalid rate limit per second value: %w", err)
		}
		cfg.RateLimitPerSecond = rate
	}

	//if crtPath, exists := os.LookupEnv("CRT_PATH"); exists {
	//	cfg.CrtPath = crtPath
	//}
	//if keyPath, exists := os.LookupEnv("KEY_PATH"); exists {
	//	cfg.KeyPath = keyPath
	//}

	return nil
}
