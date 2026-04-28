package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func LoadConfig() (*AppConfig, error) {
	loadEnvFile()

	configPath, err := resolveConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
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

func loadEnvFile() {
	candidates := []string{
		".env",
		filepath.Join("..", ".env"),
		filepath.Join("..", "..", ".env"),
	}

	for _, candidate := range candidates {
		if err := godotenv.Load(candidate); err == nil {
			return
		}
	}
}

func resolveConfigPath() (string, error) {
	candidates := []string{
		filepath.Join("Backend", "internal", "shared", "config", "DbConfig.yaml"),
		filepath.Join("..", "internal", "shared", "config", "DbConfig.yaml"),
		filepath.Join("internal", "shared", "config", "DbConfig.yaml"),
		filepath.Join("shared", "config", "DbConfig.yaml"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("DbConfig.yaml not found in expected locations")
}
