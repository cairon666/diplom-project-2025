# Миграции базы данных

- **Goose** - инструмент для управления миграциями базы данных
- **PostgreSQL** - основная база данных
- **Docker** - контейнеризация для изоляции окружения миграций
- **Docker Compose** - оркестрация контейнеров

## Работа с миграциями через Docker

### Автоматические миграции при запуске

По умолчанию миграции выполняются автоматически при запуске всего приложения:

```bash
cd deploy
docker-compose up
```

**Backend будет запущен только после успешного выполнения всех миграций.** Если миграции не удались, backend не запустится.

### Ручное управление миграций

Для ручного управления миграциями рекомендуется использовать goose напрямую (см. раздел "Работа с миграциями" ниже). Это быстрее и удобнее чем через Docker контейнеры.

## Работа с миграциями

### Установка Goose

```bash
# Установка через go install
go install github.com/pressly/goose/v3/cmd/goose@latest

# Или через brew (macOS)
brew install goose

# Или скачать бинарник с GitHub releases
```

### Настройка переменных окружения

```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
export GOOSE_MIGRATION_DIR="./backend/migrations"
```

### Создание миграции

```bash
cd backend/migrations
goose create add_users_table sql
```

### Применение миграций

```bash
cd backend/migrations

# Применить все миграции
goose up

# Применить одну миграцию
goose up-by-one

# Применить до конкретной версии
goose up-to 20240101000000
```

### Откат миграций

```bash
cd backend/migrations

# Откат одной миграции
goose down

# Откат до конкретной версии
goose down-to 20240101000000

# Сброс всех миграций
goose reset
```

### Проверка статуса

```bash
cd backend/migrations

# Статус миграций
goose status

# Версия базы данных
goose version

# Валидация
goose validate
```

### Прямое выполнение команд

```bash
cd backend/migrations

# Выполнить конкретную миграцию без версионности
goose -no-versioning up

# Исправить проблемы с версиями
goose fix
```

## Конфигурация

### Docker (docker-compose.yaml)

```yaml
environment:
  GOOSE_DRIVER: postgres
  GOOSE_MIGRATION_DIR: /migrations
  GOOSE_DBSTRING: postgresql://postgres:postgres@postgres:5432/postgres?sslmode=disable
  GOOSE_DB_HOST: postgres
  GOOSE_DB_PORT: 5432
  GOOSE_DB_USER: postgres
```

### Локальное окружение (.env)

```env
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable
GOOSE_MIGRATION_DIR=./backend/migrations
```
