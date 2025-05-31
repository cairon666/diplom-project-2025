package telegram_service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	"github.com/cairon666/vkr-backend/internal/config"
	"github.com/cairon666/vkr-backend/internal/models"
)

type TelegramService struct {
	config *config.Config
}

func NewTelegramService(config *config.Config) *TelegramService {
	return &TelegramService{config: config}
}
	
func (s *TelegramService) Verify(data models.TelegramAuthData) bool {
	if s.config.DeveloperMode {
		return true
	}

	fields := map[string]string{}

	if data.ID != 0 {
		fields["id"] = strconv.Itoa(int(data.ID))
	}

	if data.FirstName != "" {
		fields["first_name"] = data.FirstName
	}
	if data.LastName != "" {
		fields["last_name"] = data.LastName
	}
	if data.Username != "" {
		fields["username"] = data.Username
	}
	if data.PhotoURL != "" {
		fields["photo_url"] = data.PhotoURL
	}
	if data.AuthDate != 0 {
		fields["auth_date"] = strconv.Itoa(int(data.AuthDate))
	}

	// Сортировка по ключу
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Сборка строки
	var dataCheckString strings.Builder
	for i, key := range keys {
		dataCheckString.WriteString(key + "=" + fields[key])
		if i < len(keys)-1 {
			dataCheckString.WriteString("\n")
		}
	}

	// Секрет = sha256(botToken)
	secret := sha256.Sum256([]byte(s.config.Telegram.BotToken))
	h := hmac.New(sha256.New, secret[:])
	h.Write([]byte(dataCheckString.String()))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(calculatedHash), []byte(data.Hash))
}
