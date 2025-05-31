package config

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// DatabaseConfig contains database-related configuration
type DatabaseConfig struct {
	PostgresURL string
}

// ServerConfig contains web server configuration
type ServerConfig struct {
	Port int
}

// JWTConfig contains JWT-related configuration
type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
}

// LogConfig contains logging configuration
type LogConfig struct {
	Out string
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// TelegramConfig contains Telegram bot configuration
type TelegramConfig struct {
	BotToken string
}

// InfluxDBConfig contains InfluxDB configuration
type InfluxDBConfig struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

// Config represents the application configuration
type Config struct {
	Database      DatabaseConfig
	WWW           ServerConfig
	JWT           JWTConfig
	Log           LogConfig
	Redis         RedisConfig
	Telegram      TelegramConfig
	InfluxDB      InfluxDBConfig
	DeveloperMode bool
}

// Loader handles configuration loading
type Loader struct {
	k *koanf.Koanf
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	return &Loader{
		k: koanf.New("."),
	}
}

// LoadConfig loads configuration from file and environment variables
func (l *Loader) LoadConfig(filepath string) (*Config, error) {
	// Set defaults
	if err := l.setDefaults(); err != nil {
		return nil, fmt.Errorf("failed to set defaults: %w", err)
	}

	// Load from file if provided
	if filepath != "" {
		if err := l.loadFile(filepath); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Load from environment variables
	if err := l.loadFromEnv(); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Parse to config struct
	config, err := l.parseToConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate configuration
	if err := l.validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// GetDefaults returns default configuration values
func GetDefaults() map[string]interface{} {
	return map[string]interface{}{
		"postgres.url":               "postgres://user:password@localhost:5432/db",
		"www.port":                   8080,
		"jwt.secret":                 "default-secret",
		"jwt.access_token_duration":  "15m",
		"jwt.refresh_token_duration": "7d",
		"jwt.issuer":                 "default-issuer",
		"log.out":                    "console",
		"redis.addr":                 "localhost:6379",
		"redis.password":             "",
		"redis.db":                   0,
		"telegram.bot_token":         "",
		"dev":                        false,
		"influxdb.url":               "http://localhost:8086",
		"influxdb.token":             "",
		"influxdb.org":               "default-org",
		"influxdb.bucket":            "default-bucket",
	}
}

func (l *Loader) setDefaults() error {
	defaultConfig := GetDefaults()
	return l.k.Load(confmap.Provider(defaultConfig, "."), nil)
}

func (l *Loader) loadFile(filepath string) error {
	parser, err := getParser(filepath)
	if err != nil {
		return err
	}

	return l.k.Load(file.Provider(filepath), parser)
}

func (l *Loader) loadFromEnv() error {
	return l.k.Load(env.Provider("", ".", func(s string) string {
		return s
	}), nil)
}

func (l *Loader) parseToConfig() (*Config, error) {
	var cfg Config
	
	// Manual parsing to ensure compatibility
	cfg.Database.PostgresURL = l.k.String("postgres.url")
	cfg.WWW.Port = l.k.Int("www.port")
	cfg.JWT.Secret = l.k.String("jwt.secret")
	cfg.JWT.AccessTokenDuration = l.k.Duration("jwt.access_token_duration")
	cfg.JWT.RefreshTokenDuration = l.k.Duration("jwt.refresh_token_duration")
	cfg.JWT.Issuer = l.k.String("jwt.issuer")
	cfg.Log.Out = l.k.String("log.out")
	cfg.Redis.Addr = l.k.String("redis.addr")
	cfg.Redis.Password = l.k.String("redis.password")
	cfg.Redis.DB = l.k.Int("redis.db")
	cfg.Telegram.BotToken = l.k.String("telegram.bot_token")
	cfg.InfluxDB.URL = l.k.String("influxdb.url")
	cfg.InfluxDB.Token = l.k.String("influxdb.token")
	cfg.InfluxDB.Org = l.k.String("influxdb.org")
	cfg.InfluxDB.Bucket = l.k.String("influxdb.bucket")
	cfg.DeveloperMode = l.k.Bool("DEVELOPER_MODE") || l.k.Bool("dev")

	return &cfg, nil
}

func (l *Loader) validateConfig(cfg *Config) error {
	var errs []error

	// Validate required fields
	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "default-secret" {
		errs = append(errs, errors.New("JWT secret must be set and not default"))
	}

	if cfg.WWW.Port <= 0 || cfg.WWW.Port > 65535 {
		errs = append(errs, errors.New("server port must be between 1 and 65535"))
	}

	if cfg.JWT.AccessTokenDuration <= 0 {
		errs = append(errs, errors.New("JWT access token duration must be positive"))
	}

	if cfg.JWT.RefreshTokenDuration <= 0 {
		errs = append(errs, errors.New("JWT refresh token duration must be positive"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %v", errs)
	}

	return nil
}

func getParser(filepath string) (koanf.Parser, error) {
	ext := path.Ext(filepath)
	switch ext {
	case ".yaml", ".yml":
		return yaml.Parser(), nil
	case ".json":
		return json.Parser(), nil
	default:
		return nil, fmt.Errorf("unsupported config file type %s, only .yaml and .json are supported", ext)
	}
}

// Convenience function for backward compatibility
func GetConfig(filepath string) (*Config, error) {
	loader := NewLoader()
	return loader.LoadConfig(filepath)
}

// Helper methods for Config struct
func (c *Config) IsProduction() bool {
	return !c.DeveloperMode
}

func (c *Config) GetPostgresURL() string {
	return c.Database.PostgresURL
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%d", c.WWW.Port)
}

func (c *Config) GetRedisAddress() string {
	return c.Redis.Addr
}
