FROM golang:1.23-alpine AS builder

# Устанавливаем goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates postgresql-client

# Копируем goose binary из builder stage
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Создаем рабочую директорию
WORKDIR /migrations

# Копируем файлы миграций
COPY *.sql ./

# Создаем entrypoint script
RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo 'set -e' >> /entrypoint.sh && \
    echo 'echo "Waiting for database..."' >> /entrypoint.sh && \
    echo 'until pg_isready -h "$GOOSE_DB_HOST" -p "$GOOSE_DB_PORT" -U "$GOOSE_DB_USER"; do' >> /entrypoint.sh && \
    echo '  echo "Database is unavailable - sleeping"' >> /entrypoint.sh && \
    echo '  sleep 1' >> /entrypoint.sh && \
    echo 'done' >> /entrypoint.sh && \
    echo 'echo "Database is up - executing migrations"' >> /entrypoint.sh && \
    echo 'goose "$@"' >> /entrypoint.sh && \
    echo 'MIGRATION_EXIT_CODE=$?' >> /entrypoint.sh && \
    echo 'if [ $MIGRATION_EXIT_CODE -eq 0 ]; then' >> /entrypoint.sh && \
    echo '  echo "Migrations completed successfully"' >> /entrypoint.sh && \
    echo '  # Создаем файл-маркер успешного завершения' >> /entrypoint.sh && \
    echo '  touch /tmp/migrations_completed' >> /entrypoint.sh && \
    echo '  # Держим контейнер живым для healthcheck' >> /entrypoint.sh && \
    echo '  while true; do sleep 30; done' >> /entrypoint.sh && \
    echo 'else' >> /entrypoint.sh && \
    echo '  echo "Migrations failed with exit code $MIGRATION_EXIT_CODE"' >> /entrypoint.sh && \
    echo '  exit $MIGRATION_EXIT_CODE' >> /entrypoint.sh && \
    echo 'fi' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

# Создаем healthcheck script
RUN echo '#!/bin/sh' > /healthcheck.sh && \
    echo 'if [ -f /tmp/migrations_completed ]; then' >> /healthcheck.sh && \
    echo '  echo "Migrations completed successfully"' >> /healthcheck.sh && \
    echo '  exit 0' >> /healthcheck.sh && \
    echo 'else' >> /healthcheck.sh && \
    echo '  echo "Migrations not completed yet"' >> /healthcheck.sh && \
    echo '  exit 1' >> /healthcheck.sh && \
    echo 'fi' >> /healthcheck.sh && \
    chmod +x /healthcheck.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["up"] 