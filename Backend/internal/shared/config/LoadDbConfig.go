package config

import (
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()

	data, err := os.ReadFile("shared/config/DbConfig.yaml")
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.Database.Password = os.Getenv("DB_PASSWORD")

	return &cfg, nil
}
