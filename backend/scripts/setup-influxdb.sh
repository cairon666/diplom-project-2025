#!/bin/bash

# Скрипт для настройки InfluxDB после запуска docker-compose

set -e

echo "🚀 Настройка InfluxDB для разработки..."

# Ждем пока InfluxDB запустится
echo "⏳ Ожидание запуска InfluxDB..."
until curl -f http://localhost:8086/ping > /dev/null 2>&1; do
    echo "   InfluxDB еще не готов, ждем..."
    sleep 2
done

echo "✅ InfluxDB запущен!"

# Проверяем, что организация и bucket созданы
echo "🔍 Проверка конфигурации..."

# Устанавливаем переменные (соответствуют config/config.local.yaml)
INFLUX_URL="http://localhost:8086"
INFLUX_TOKEN="dev-token-for-local-development-only"
INFLUX_ORG="health-analytics"
INFLUX_BUCKET="health-data"

# Проверяем bucket
if docker exec influxdb influx bucket list --org $INFLUX_ORG --token $INFLUX_TOKEN | grep -q $INFLUX_BUCKET; then
    echo "✅ Bucket '$INFLUX_BUCKET' уже существует"
else
    echo "📦 Создание bucket '$INFLUX_BUCKET'..."
    docker exec influxdb influx bucket create \
        --name $INFLUX_BUCKET \
        --org $INFLUX_ORG \
        --token $INFLUX_TOKEN \
        --retention 8760h
    echo "✅ Bucket создан"
fi
