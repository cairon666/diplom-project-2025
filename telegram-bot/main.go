package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cairon666/vkr-backend/telegram-bot/internal/bot"
	"github.com/cairon666/vkr-backend/telegram-bot/internal/config"
)

func main() {
	// Парсим флаги командной строки
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Настраиваем логгер
	logger, err := setupLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		log.Fatalf("Failed to setup logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Telegram Bot",
		zap.String("config_path", *configPath),
		zap.String("log_level", cfg.Log.Level),
	)

	// Создаем бота
	telegramBot, err := bot.NewBot(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to create bot", zap.Error(err))
	}

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем бота в отдельной горутине
	go func() {
		if err := telegramBot.Start(); err != nil {
			logger.Error("Bot error", zap.Error(err))
			cancel()
		}
	}()

	// Ожидаем сигнал завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logger.Info("Received signal, shutting down", zap.String("signal", sig.String()))
	case <-ctx.Done():
		logger.Info("Context cancelled, shutting down")
	}

	// Останавливаем бота
	telegramBot.Stop()
	logger.Info("Bot stopped successfully")
}

// setupLogger настраивает логгер
func setupLogger(level, format string) (*zap.Logger, error) {
	// Парсим уровень логирования
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Настраиваем конфигурацию логгера
	var config zap.Config
	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// Создаем логгер
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
} 