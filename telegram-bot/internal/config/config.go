package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	API      APIConfig      `yaml:"api"`
	Log      LogConfig      `yaml:"log"`
	Storage  StorageConfig  `yaml:"storage"`
}

type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	Debug    bool   `yaml:"debug"`
}

type APIConfig struct {
	BaseURL string        `yaml:"base_url"`
	Timeout time.Duration `yaml:"timeout"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type StorageConfig struct {
	FilePath string `yaml:"file_path"`
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Проверяем обязательные поля
	if config.Telegram.BotToken == "" || config.Telegram.BotToken == "YOUR_BOT_TOKEN_HERE" {
		return nil, fmt.Errorf("telegram bot token is required")
	}

	if config.API.BaseURL == "" {
		return nil, fmt.Errorf("api base url is required")
	}

	// Устанавливаем значения по умолчанию
	if config.API.Timeout == 0 {
		config.API.Timeout = 30 * time.Second
	}

	if config.Log.Level == "" {
		config.Log.Level = "info"
	}

	if config.Storage.FilePath == "" {
		config.Storage.FilePath = "./data/users.json"
	}

	return &config, nil
} 