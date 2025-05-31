package bot

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/cairon666/vkr-backend/telegram-bot/internal/api"
	"github.com/cairon666/vkr-backend/telegram-bot/internal/config"
	"github.com/cairon666/vkr-backend/telegram-bot/internal/storage"
)

// Bot представляет Telegram бота
type Bot struct {
	api     *tgbotapi.BotAPI
	client  *api.Client
	storage *storage.Storage
	logger  *zap.Logger
	config  *config.Config
}

// NewBot создает новый экземпляр бота
func NewBot(cfg *config.Config, logger *zap.Logger) (*Bot, error) {
	// Создаем Telegram API клиент
	botAPI, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	botAPI.Debug = cfg.Telegram.Debug

	// Создаем HTTP клиент для API
	apiClient := api.NewClient(cfg.API.BaseURL, cfg.API.Timeout)

	// Создаем хранилище
	store, err := storage.NewStorage(cfg.Storage.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	return &Bot{
		api:     botAPI,
		client:  apiClient,
		storage: store,
		logger:  logger,
		config:  cfg,
	}, nil
}

// Start запускает бота
func (b *Bot) Start() error {
	b.logger.Info("Starting bot", zap.String("username", b.api.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}

	return nil
}

// Stop останавливает бота
func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
	b.logger.Info("Bot stopped")
}

// handleMessage обрабатывает входящие сообщения
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	userID := message.From.ID
	chatID := message.Chat.ID

	b.logger.Info("Received message",
		zap.Int64("user_id", userID),
		zap.Int64("chat_id", chatID),
		zap.String("text", message.Text),
	)

	// Обрабатываем команды
	if message.IsCommand() {
		b.handleCommand(message)
		return
	}

	// Проверяем, есть ли у пользователя API ключ
	user, exists := b.storage.GetUser(userID)
	if !exists {
		b.sendMessage(chatID, "❌ Вы не авторизованы. Используйте команду /start для начала работы.")
		return
	}

	if user.APIKey == "" {
		b.sendMessage(chatID, "❌ API ключ не установлен. Используйте команду /connect для подключения.")
		return
	}

	// Обрабатываем обычные сообщения
	b.handleRegularMessage(message, user)
}

// handleCommand обрабатывает команды
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := message.From.ID

	switch message.Command() {
	case "start":
		b.handleStartCommand(chatID, userID, message.From)
	case "connect":
		b.handleConnectCommand(chatID, userID, message.CommandArguments())
	case "disconnect":
		b.handleDisconnectCommand(chatID, userID)
	case "status":
		b.handleStatusCommand(chatID, userID)
	case "help":
		b.handleHelpCommand(chatID)
	case "rrintervals", "rr":
		b.handleRRIntervalsCommand(chatID, userID, message.CommandArguments())
	case "rrstats", "rrstatistics":
		b.handleRRStatisticsCommand(chatID, userID, message.CommandArguments())
	case "rranalyze", "rranalysis":
		b.handleRRAnalysisCommand(chatID, userID, message.CommandArguments())
	default:
		b.sendMessage(chatID, "❌ Неизвестная команда. Используйте /help для просмотра доступных команд.")
	}
}

// handleStartCommand обрабатывает команду /start
func (b *Bot) handleStartCommand(chatID, userID int64, from *tgbotapi.User) {
	// Сохраняем пользователя в хранилище
	user := &storage.UserData{
		TelegramID: userID,
		Username:   from.UserName,
		FirstName:  from.FirstName,
		LastName:   from.LastName,
		CreatedAt:  time.Now().Format(time.RFC3339),
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}

	if err := b.storage.SaveUser(user); err != nil {
		b.logger.Error("Failed to save user", zap.Error(err))
		b.sendMessage(chatID, "❌ Произошла ошибка при сохранении данных.")
		return
	}

	message := `�� Добро пожаловать в R-R Intervals Analytics Bot!

Этот бот специализируется на анализе R-R интервалов и вариабельности сердечного ритма (HRV).

Для начала работы вам необходимо:

1️⃣ Авторизоваться на веб-сайте:
🔗 https://your-frontend-url.com/login

2️⃣ Создать External App и получить API ключ:
🔗 https://your-frontend-url.com/external-apps

3️⃣ Подключить бота с помощью команды:
/connect YOUR_API_KEY

После подключения вы сможете анализировать R-R интервалы прямо в Telegram!

Используйте /help для просмотра всех доступных команд.`

	b.sendMessage(chatID, message)
}

// handleConnectCommand обрабатывает команду /connect
func (b *Bot) handleConnectCommand(chatID, userID int64, args string) {
	if args == "" {
		b.sendMessage(chatID, "❌ Укажите API ключ: /connect YOUR_API_KEY")
		return
	}

	apiKey := strings.TrimSpace(args)

	// Проверяем формат API ключа (base64-подобная строка)
	apiKeyRegex := regexp.MustCompile(`^[A-Za-z0-9+/\-_]+=*$`)
	if len(apiKey) < 20 || len(apiKey) > 100 || !apiKeyRegex.MatchString(apiKey) {
		b.sendMessage(chatID, "❌ Неверный формат API ключа. Ключ должен содержать буквы, цифры и символы +/-/=/_.")
		return
	}

	// Проверяем, что API ключ работает (делаем тестовый запрос)
	now := time.Now()
	from := now.AddDate(0, 0, -1) // вчера
	to := now

	// Попробуем сделать простой запрос для проверки ключа
	_, err := b.client.GetRRStatistics(apiKey, from, to, false, false, 0)
	if err != nil {
		// Логируем для отладки, но не блокируем если это просто отсутствие данных
		b.logger.Warn("API key validation request failed", zap.Error(err))
		
		// Проверяем, что это не критическая ошибка авторизации
		errStr := err.Error()
		if strings.Contains(errStr, "status 401") || strings.Contains(errStr, "status 403") || strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "forbidden") {
			b.logger.Error("Failed to validate API key - unauthorized", zap.Error(err))
			b.sendMessage(chatID, "❌ Неверный API ключ.")
			return
		}
		
		// Если это ошибка формата или отсутствия данных - продолжаем
		if strings.Contains(errStr, "format") || strings.Contains(errStr, "status 404") || strings.Contains(errStr, "no data") {
			b.logger.Info("API key validation - format issue or no data, but key seems valid", zap.Error(err))
		}
	}

	// Сохраняем API ключ
	user, exists := b.storage.GetUser(userID)
	if !exists {
		b.sendMessage(chatID, "❌ Пользователь не найден. Используйте /start для регистрации.")
		return
	}

	user.APIKey = apiKey
	user.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := b.storage.SaveUser(user); err != nil {
		b.logger.Error("Failed to save user with API key", zap.Error(err))
		b.sendMessage(chatID, "❌ Произошла ошибка при сохранении API ключа.")
		return
	}

	b.sendMessage(chatID, "✅ API ключ успешно подключен! Теперь вы можете анализировать R-R интервалы.")
}

// handleDisconnectCommand обрабатывает команду /disconnect
func (b *Bot) handleDisconnectCommand(chatID, userID int64) {
	user, exists := b.storage.GetUser(userID)
	if !exists {
		b.sendMessage(chatID, "❌ Пользователь не найден.")
		return
	}

	user.APIKey = ""
	user.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := b.storage.SaveUser(user); err != nil {
		b.logger.Error("Failed to disconnect user", zap.Error(err))
		b.sendMessage(chatID, "❌ Произошла ошибка при отключении.")
		return
	}

	b.sendMessage(chatID, "✅ API ключ отключен. Используйте /connect для повторного подключения.")
}

// handleStatusCommand обрабатывает команду /status
func (b *Bot) handleStatusCommand(chatID, userID int64) {
	user, exists := b.storage.GetUser(userID)
	if !exists {
		b.sendMessage(chatID, "❌ Пользователь не найден. Используйте /start для регистрации.")
		return
	}

	status := "❌ Не подключен"
	if user.APIKey != "" {
		status = "✅ Подключен"
	}

	message := fmt.Sprintf(`📊 Статус подключения: %s

👤 Пользователь: %s %s
🆔 Telegram ID: %d
📅 Зарегистрирован: %s
🔄 Обновлен: %s`,
		status,
		user.FirstName, user.LastName,
		user.TelegramID,
		user.CreatedAt,
		user.UpdatedAt,
	)

	b.sendMessage(chatID, message)
}

// handleHelpCommand обрабатывает команду /help
func (b *Bot) handleHelpCommand(chatID int64) {
	message := `📋 Доступные команды:

🔧 Управление:
/start - Начать работу с ботом
/connect API_KEY - Подключить API ключ
/disconnect - Отключить API ключ
/status - Проверить статус подключения
/help - Показать эту справку

💓 R-R интервалы:
/rr [дата_начала] [дата_конца] - Получить R-R интервалы
/rrstats [дата_начала] [дата_конца] - Статистика R-R интервалов  
/rranalyze [дата_начала] [дата_конца] - Полный анализ R-R интервалов

Примеры:
/rr 2024-01-01 2024-01-02
/rrstats 2024-01-01
/rranalyze 2024-01-01 2024-01-02

Форматы дат:
- Без аргументов: сегодня
- Одна дата: весь день (2024-01-01)
- Две даты: диапазон (2024-01-01 2024-01-02)`

	b.sendMessage(chatID, message)
}

// handleRegularMessage обрабатывает обычные сообщения
func (b *Bot) handleRegularMessage(message *tgbotapi.Message, user *storage.UserData) {
	chatID := message.Chat.ID
	text := strings.ToLower(strings.TrimSpace(message.Text))

	switch {
	case strings.Contains(text, "r-r") || strings.Contains(text, "rr интервал") || strings.Contains(text, "вариабельность") || strings.Contains(text, "hrv"):
		b.sendMessage(chatID, `💓 *Работа с R-R интервалами:*

/rr - получить сырые данные
/rrstats - базовая статистика и HRV
/rranalyze - полный анализ с интерпретацией

Примеры:
/rr 2024-01-01 2024-01-02
/rrstats 2024-01-01
/rranalyze вчера сегодня`)
	case strings.Contains(text, "помощь") || strings.Contains(text, "help") || strings.Contains(text, "команды"):
		b.handleHelpCommand(chatID)
	default:
		b.sendMessage(chatID, `❓ Не понимаю. 

Этот бот работает только с R-R интервалами. 
Используйте /help для просмотра команд или напишите "r-r интервалы" для справки.`)
	}
}

// sendMessage отправляет сообщение пользователю
func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := b.api.Send(msg); err != nil {
		b.logger.Error("Failed to send message", zap.Error(err))
	}
}
