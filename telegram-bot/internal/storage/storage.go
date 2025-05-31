package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// UserData представляет данные пользователя
type UserData struct {
	TelegramID int64  `json:"telegram_id"`
	APIKey     string `json:"api_key"`
	Username   string `json:"username,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// Storage представляет файловое хранилище
type Storage struct {
	filePath string
	mu       sync.RWMutex
	users    map[int64]*UserData
}

// NewStorage создает новое файловое хранилище
func NewStorage(filePath string) (*Storage, error) {
	// Создаем директорию если не существует
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	storage := &Storage{
		filePath: filePath,
		users:    make(map[int64]*UserData),
	}

	// Загружаем существующие данные
	if err := storage.load(); err != nil {
		return nil, fmt.Errorf("failed to load storage: %w", err)
	}

	return storage, nil
}

// SaveUser сохраняет данные пользователя
func (s *Storage) SaveUser(user *UserData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.TelegramID] = user
	return s.save()
}

// GetUser получает данные пользователя по Telegram ID
func (s *Storage) GetUser(telegramID int64) (*UserData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[telegramID]
	return user, exists
}

// GetUserByAPIKey получает данные пользователя по API ключу
func (s *Storage) GetUserByAPIKey(apiKey string) (*UserData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.APIKey == apiKey {
			return user, true
		}
	}
	return nil, false
}

// DeleteUser удаляет пользователя
func (s *Storage) DeleteUser(telegramID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.users, telegramID)
	return s.save()
}

// GetAllUsers возвращает всех пользователей
func (s *Storage) GetAllUsers() []*UserData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*UserData, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// load загружает данные из файла
func (s *Storage) load() error {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// Файл не существует, создаем пустое хранилище
		return nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		// Пустой файл
		return nil
	}

	var users []*UserData
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	// Преобразуем в map
	for _, user := range users {
		s.users[user.TelegramID] = user
	}

	return nil
}

// save сохраняет данные в файл
func (s *Storage) save() error {
	// Преобразуем map в slice
	users := make([]*UserData, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
} 