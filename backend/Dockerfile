# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.24-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Создаем пользователя для безопасности
RUN adduser -D -g '' appuser

# Устанавливаем рабочую директорию
WORKDIR /build

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app ./cmd/www/main.go

# Финальный образ
FROM scratch

# Копируем сертификаты для HTTPS запросов
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем информацию о временных зонах
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Копируем пользователя
COPY --from=builder /etc/passwd /etc/passwd

# Копируем скомпилированное приложение
COPY --from=builder /build/app /app

# Создаем директорию для конфигурации
USER appuser

# Открываем порт
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["/app"] 