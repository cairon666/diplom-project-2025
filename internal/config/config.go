package config

import (
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	PostgresURL string
	WWW         struct {
		Port int
	}
	JWT struct {
		Secret               string
		AccessTokenDuration  time.Duration
		RefreshTokenDuration time.Duration
		Issuer               string
	}
	Log struct {
		Out string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	Telegram struct {
		BotToken string
	}
	DeveloperMode bool
}

var (
	k         = koanf.New(".")
	config    *Config
	globalErr error
	once      sync.Once
)

func GetConfig(filepath string) (*Config, error) {
	once.Do(func() {
		setDefaults()
		if err := loadFile(filepath); err != nil {
			globalErr = err
			return
		}

		if err := k.Load(env.Provider("", ".", func(s string) string {
			return s
		}), nil); err != nil {
			globalErr = err
			return
		}

		config, globalErr = parseToConfig(k)
	})

	return config, globalErr
}

func setDefaults() {
	defaultConfig := map[string]interface{}{
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
	}

	k.Load(confmap.Provider(defaultConfig, "."), nil)
}

func loadFile(filepath string) error {
	parser, err := getParser(filepath)
	if err != nil {
		return err
	}

	if err := k.Load(file.Provider(filepath), parser); err != nil {
		return err
	}

	return nil
}

func parseToConfig(k *koanf.Koanf) (*Config, error) {
	var cfg Config

	// postgres
	cfg.PostgresURL = k.String("postgres.url")
	// www
	cfg.WWW.Port = k.Int("www.port")
	// jwt
	cfg.JWT.Secret = k.String("jwt.secret")
	cfg.JWT.AccessTokenDuration = k.Duration("jwt.access_token_duration")
	cfg.JWT.RefreshTokenDuration = k.Duration("jwt.refresh_token_duration")
	cfg.JWT.Issuer = k.String("jwt.issuer")
	// log
	cfg.Log.Out = k.String("log.out")
	// redis
	cfg.Redis.Addr = k.String("redis.addr")
	cfg.Redis.Password = k.String("redis.password")
	cfg.Redis.DB = k.Int("redis.db")
	// telegram
	cfg.Telegram.BotToken = k.String("telegram.bot_token")
	// dev
	cfg.DeveloperMode = k.Bool("DEVELOPER_MODE") || k.Bool("dev")

	return &cfg, nil
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
