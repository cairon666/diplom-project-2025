# CI/CD Pipeline

Этот репозиторий использует GitHub Actions для автоматизации проверки качества кода.

## 🔄 CI Pipeline (`ci.yml`)

Запускается при:

- Push в ветки `main`, `develop`
- Создании Pull Request'а

**Включает два параллельных job'а:**

### 🔧 Backend

- 📦 **Установка зависимостей**: `go mod download`
- 🧹 **Линтинг**: golangci-lint с настройками из `.golangci.yaml`
- 🧪 **Тестирование**: `go test -v ./...`
- 🏗️ **Сборка**: `go build -v ./cmd/...`

### ⚛️ Frontend

- 📦 **Установка зависимостей**: `npm ci`
- 🔍 **TypeScript проверка**: `npm run tsc`
- 🧹 **Линтинг**: ESLint
- 🧪 **Тестирование**: `npm run test`
- 🏗️ **Сборка**: `npm run build`

## 🛡️ Защита веток

### Настройка обязательна!

Следуйте инструкциям в [BRANCH_PROTECTION_SETUP.md](./BRANCH_PROTECTION_SETUP.md) для настройки защиты веток.

**Обязательные статус-чеки:**

- ✅ Backend
- ✅ Frontend

## 🚀 Локальная разработка

### Подготовка окружения

```bash
# Backend
cd backend
go mod download
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Frontend
cd frontend
npm ci
```

### Запуск проверок локально

```bash
# Backend
cd backend
golangci-lint run
go test ./...
go build -v ./cmd/...

# Frontend
cd frontend
npm run tsc
npx eslint . --ext .ts,.tsx
npm run test
npm run build
```

### Pre-commit проверки

Используйте lefthook для автоматических проверок:

```bash
# Установка (если еще не установлен)
go install github.com/evilmartians/lefthook@latest
lefthook install

# Теперь проверки будут запускаться автоматически при commit
```

## 🔧 Настройка для нового проекта

1. **Обновите CODEOWNERS**: Замените `@cairon666` на ваших пользователей
2. **Настройте защиту веток**: Следуйте инструкциям в BRANCH_PROTECTION_SETUP.md

## 📝 Соглашения

### Branch naming

- `feature/TASK-123-description` - новые функции
- `bugfix/TASK-456-description` - исправления багов
- `hotfix/critical-issue` - критические исправления

### Commit messages

Используйте conventional commits:

- `feat: add user authentication`
- `fix: resolve memory leak in backend`
- `docs: update API documentation`
- `test: add integration tests for orders`

### Pull Requests

- Используйте понятные заголовки
- Заполняйте описание с контекстом изменений
- Добавляйте скриншоты для UI изменений
- Привязывайте к задачам (issue/task)

## 🆘 Troubleshooting

### CI падает с ошибкой зависимостей

1. Проверьте версии в `go.mod` и `package.json`
2. Убедитесь что `package-lock.json` актуален
3. Локально запустите `go mod tidy` и `npm ci`

### Тесты проходят локально, но падают в CI

1. Убедитесь что тесты не зависят от локального состояния
2. Проверьте пути к файлам (case sensitivity)
3. Проверьте временные зоны и локали

### Линтер находит ошибки

1. Запустите линтер локально: `golangci-lint run` или `npm run lint`
2. Исправьте найденные проблемы
3. Убедитесь что настройки линтера одинаковые

## 📞 Контакты

По вопросам CI/CD обращайтесь к @cairon666
