# Logger Package

Улучшенная архитектура логгера с поддержкой интерфейсов для лучшей тестируемости.

## Преимущества новой архитектуры

- ✅ **Интерфейс-ориентированный дизайн** - легко тестировать и подменять реализации
- ✅ **Mock логгер** с возможностью проверки вызовов
- ✅ **No-op логгер** для быстрых тестов
- ✅ **Поддержка контекста** для tracing и correlation
- ✅ **Структурированное логирование** с типизированными полями
- ✅ **HTTP middleware** для автоматического логирования запросов

## Файлы

- `interface.go` - Основные интерфейсы логгера
- `logger.go` - Основная реализация на базе Zap
- `mock.go` - Mock логгер для тестирования
- `integration.go` - Интеграция с HTTP и контекстом
- `fields.go` - Вспомогательные функции для полей
- `examples.go` - Примеры использования

## Основные интерфейсы

### ILogger

Базовый интерфейс для всех логгеров:

```go
type ILogger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Debug(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    With(fields ...Field) ILogger
    WithContext(ctx context.Context) ILogger
}
```

### FieldLogger

Расширенный интерфейс с дополнительными уровнями:

```go
type FieldLogger interface {
    ILogger
    Fatal(msg string, fields ...Field)
    Panic(msg string, fields ...Field)
}
```

### ContextLogger

Интерфейс для работы с контекстом:

```go
type ContextLogger interface {
    InfoContext(ctx context.Context, msg string, fields ...Field)
    ErrorContext(ctx context.Context, msg string, fields ...Field)
    // ...
}
```

## Использование

### 1. В продакшн коде

```go
// Старый способ ❌
func NewUserService(db *sql.DB, logger *logger.Logger) *UserService {
    return &UserService{db: db, logger: logger}
}

// Новый способ ✅
func NewUserService(db *sql.DB, logger logger.ILogger) *UserService {
    return &UserService{db: db, logger: logger}
}

func main() {
    // Создаем настоящий логгер
    logger, err := logger.NewProd()
    if err != nil {
        panic(err)
    }

    // Используем через интерфейс
    userService := NewUserService(db, logger)
}
```

### 2. В тестах

```go
func TestUserService_GetByID(t *testing.T) {
    // Создаем mock логгер
    mockLogger := logger.NewMockLogger()

    // Создаем сервис с mock логгером
    service := NewUserService(mockDB, mockLogger)

    // Тестируем
    user, err := service.GetByID(123)

    // Проверяем логи
    assert.Equal(t, 1, mockLogger.GetCallsCount())
    assert.True(t, mockLogger.HasCallWithMessage("Getting user by ID"))

    calls := mockLogger.GetCallsByLevel("INFO")
    assert.Len(t, calls, 1)
}
```

### 3. Структурированное логирование

```go
func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
    log := logger.WithContextFields(s.logger, ctx)

    log.Info("Getting user by ID", logger.WithUserID(id))

    user, err := s.repo.GetByID(id)
    if err != nil {
        log.Error("Failed to get user",
            logger.WithError(err),
            logger.WithUserID(id),
        )
        return nil, apperrors.InternalErrorf("Failed to get user: %v", err)
    }

    log.Info("User retrieved successfully", logger.WithUserID(id))
    return user, nil
}
```

### 4. HTTP Middleware

```go
func main() {
    logger, _ := logger.NewProd()

    r := gin.New()
    r.Use(logger.GinMiddleware(logger))

    r.GET("/users/:id", getUserHandler)
    r.Run()
}
```

### 5. LoggedService

```go
type UserService struct {
    *logger.LoggedService
    repo UserRepository
}

func NewUserService(repo UserRepository, logger logger.ILogger) *UserService {
    return &UserService{
        LoggedService: logger.NewLoggedService(logger, "user_service"),
        repo:          repo,
    }
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    s.LogInfo(ctx, "Creating new user", logger.String("email", user.Email))

    if err := s.repo.Create(user); err != nil {
        s.LogError(ctx, err, "Failed to create user",
            logger.String("email", user.Email),
        )
        return err
    }

    s.LogInfo(ctx, "User created successfully",
        logger.String("email", user.Email),
        logger.WithUserID(user.ID),
    )
    return nil
}
```

## Доступные типы логгеров

### 1. Production Logger

```go
logger, err := logger.NewProd()  // JSON, файлы, INFO+ уровень
```

### 2. Development Logger

```go
logger, err := logger.NewDev()   // Console, stderr, DEBUG+ уровень
```

### 3. Mock Logger

```go
mockLogger := logger.NewMockLogger()  // Для тестирования
```

### 4. No-Op Logger

```go
noopLogger := &logger.NoOpLogger{}    // Заглушка, ничего не делает
```

## Вспомогательные функции для полей

```go
logger.String("key", "value")
logger.Int("count", 42)
logger.WithError(err)
logger.WithUserID(123)
logger.WithRequestID("req-123")
logger.WithOperation("create_user")
logger.WithDuration(time.Since(start))
```

## Интеграция с контекстом

```go
// Извлечение полей из контекста
fields := logger.FromContext(ctx)

// Создание логгера с контекстными полями
log := logger.WithContextFields(baseLogger, ctx)

// Использование в middleware
ctx := context.WithValue(ctx, logger.RequestIDKey, "req-123")
ctx = context.WithValue(ctx, logger.UserIDKey, int64(456))
```

## Mock Logger API

```go
mockLogger := logger.NewMockLogger()

// Получение всех вызовов
calls := mockLogger.GetCalls()

// Получение последнего вызова
lastCall := mockLogger.GetLastCall()

// Количество вызовов
count := mockLogger.GetCallsCount()

// Вызовы по уровню
infoCalls := mockLogger.GetCallsByLevel("INFO")

// Поиск по сообщению
hasCall := mockLogger.HasCallWithMessage("User created")

// Очистка истории
mockLogger.Clear()
```

## Миграция

### Шаг 1: Замените типы на интерфейсы

```go
// Было
func NewService(logger *logger.Logger) *Service

// Стало
func NewService(logger logger.ILogger) *Service
```

### Шаг 2: Обновите тесты

```go
// Было
logger, _ := logger.NewDev()
service := NewService(logger)

// Стало
mockLogger := logger.NewMockLogger()
service := NewService(mockLogger)
```

### Шаг 3: Добавьте контекст

```go
// Было
func (s *Service) DoSomething(id int) error

// Стало
func (s *Service) DoSomething(ctx context.Context, id int) error
```

## Лучшие практики

1. **Всегда используйте интерфейсы** вместо конкретных типов
2. **Передавайте контекст** для correlation и tracing
3. **Используйте структурированные поля** вместо форматированных строк
4. **Логируйте на правильном уровне**:
   - DEBUG - детальная информация для отладки
   - INFO - общая информация о работе системы
   - WARN - потенциальные проблемы
   - ERROR - ошибки, которые нужно исправить
5. **В тестах проверяйте важные логи** для критичной логики
6. **Используйте No-Op логгер** когда логи не важны для теста
