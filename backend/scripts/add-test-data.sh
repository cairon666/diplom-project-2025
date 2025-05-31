#!/bin/bash

# Скрипт для добавления тестовых R-R интервалов в InfluxDB

set -e

echo "💓 Добавление R-R интервалов в InfluxDB..."

# Переменные (соответствуют config/config.local.yaml)
INFLUX_URL="http://localhost:8086"
INFLUX_TOKEN="dev-token-for-local-development-only"
INFLUX_ORG="health-analytics"
INFLUX_BUCKET="health-data"

# Тестовые UUID
USER_ID="b14adb16-5a49-4d8e-a97b-e006415ec126"
DEVICE_ID="fa9e1b6d-a593-469e-b174-f5c81028d3c9"

echo "👤 Пользователь ID: $USER_ID"
echo "📱 Устройство ID: $DEVICE_ID"
echo ""

# Определяем ОС для совместимости команд date
OS_TYPE=$(uname -s)

# Функция для работы с датами и временем в зависимости от ОС
date_with_time() {
    local days=$1
    local hour=$2
    local minute=${3:-0}
    local second=${4:-0}
    
    if [ "$OS_TYPE" = "Darwin" ]; then
        # macOS
        local base_date=$(date -v-${days}d '+%Y-%m-%d')
        date -j -f '%Y-%m-%d %H:%M:%S' "${base_date} ${hour}:${minute}:${second}" '+%s'
    else
        # Linux
        date -d "${days} days ago ${hour}:${minute}:${second}" '+%s'
    fi
}

# Функция для генерации UUID v4
generate_uuid() {
    if command -v python3 &> /dev/null; then
        python3 -c "import uuid; print(str(uuid.uuid4()))"
    elif command -v uuidgen &> /dev/null; then
        uuidgen | tr '[:upper:]' '[:lower:]'
    else
        # Fallback для систем без python3 и uuidgen
        echo "$(od -x /dev/urandom | head -1 | awk '{OFS="-"; print $2$3,$4,$5,$6,$7$8$9}')"
    fi
}

# Функция для генерации реалистичного R-R интервала в зависимости от времени и активности
generate_rr_interval() {
    local hour=$1
    local base_rr
    local variation
    
    # Определяем базовый R-R интервал в зависимости от времени суток
    if [ $hour -ge 0 ] && [ $hour -le 6 ]; then
        # Ночь/раннее утро: медленный пульс, длинные R-R интервалы
        base_rr=1100  # ~55 BPM
        variation=200
    elif [ $hour -ge 7 ] && [ $hour -le 9 ]; then
        # Утреннее пробуждение: умеренный пульс
        base_rr=900   # ~67 BPM
        variation=150
    elif [ $hour -ge 10 ] && [ $hour -le 12 ]; then
        # Утренняя активность: немного быстрее
        base_rr=800   # ~75 BPM
        variation=100
    elif [ $hour -ge 13 ] && [ $hour -le 17 ]; then
        # Дневная активность: может быть физическая нагрузка
        base_rr=750   # ~80 BPM
        variation=150
    elif [ $hour -ge 18 ] && [ $hour -le 21 ]; then
        # Вечер: умеренный
        base_rr=850   # ~71 BPM
        variation=120
    else
        # Поздний вечер: подготовка ко сну
        base_rr=950   # ~63 BPM
        variation=100
    fi
    
    # Добавляем случайную вариацию
    local random_offset=$((RANDOM % (variation * 2) - variation))
    local rr_interval=$((base_rr + random_offset))
    
    # Иногда добавляем "события" для реалистичности
    local event_chance=$((RANDOM % 100))
    if [ $event_chance -lt 5 ]; then
        # 5% шанс на экстрасистолу (короткий интервал)
        rr_interval=$((rr_interval - 200 - RANDOM % 300))
    elif [ $event_chance -lt 10 ]; then
        # 5% шанс на компенсаторную паузу (длинный интервал)
        rr_interval=$((rr_interval + 300 + RANDOM % 400))
    elif [ $event_chance -lt 20 ]; then
        # 10% шанс на стресс/физическую активность (короткие интервалы)
        rr_interval=$((rr_interval - 150 - RANDOM % 200))
    fi
    
    # Ограничиваем физиологическими пределами (300-2000 мс)
    if [ $rr_interval -lt 300 ]; then
        rr_interval=300
    elif [ $rr_interval -gt 2000 ]; then
        rr_interval=2000
    fi
    
    echo $rr_interval
}

# Функция для генерации качества сигнала
generate_signal_quality() {
    local hour=$1
    local base_quality
    
    # Качество сигнала зависит от времени (ночью лучше, днем хуже из-за движения)
    if [ $hour -ge 0 ] && [ $hour -le 6 ]; then
        base_quality=0.9  # Высокое качество ночью
    elif [ $hour -ge 7 ] && [ $hour -le 9 ]; then
        base_quality=0.8  # Хорошее качество утром
    elif [ $hour -ge 10 ] && [ $hour -le 17 ]; then
        base_quality=0.7  # Среднее качество днем (движение)
    else
        base_quality=0.85 # Хорошее качество вечером
    fi
    
    # Добавляем случайную вариацию ±0.15
    local variation=$((RANDOM % 30 - 15))  # -15 to +15
    local quality=$(echo "scale=2; $base_quality + $variation / 100" | bc -l 2>/dev/null || echo $base_quality)
    
    # Ограничиваем диапазоном 0.3-1.0
    if (( $(echo "$quality < 0.3" | bc -l 2>/dev/null || echo 0) )); then
        quality=0.3
    elif (( $(echo "$quality > 1.0" | bc -l 2>/dev/null || echo 0) )); then
        quality=1.0
    fi
    
    echo $quality
}

echo "💓 Добавление R-R интервалов..."
TOTAL_RECORDS=0
VALID_INTERVALS=0
ANOMALOUS_INTERVALS=0
HIGH_QUALITY_COUNT=0

# Генерируем данные за последние 30 дней
for i in {0..29}; do
    echo "📅 Обрабатываем день -$i..."
    
    for hour in {0..23}; do  # Данные круглосуточно
        # Генерируем 50-120 R-R интервалов в час (имитация непрерывного мониторинга)
        intervals_per_hour=$((50 + RANDOM % 70))
        
        for ((j=0; j<intervals_per_hour; j++)); do
            # Генерируем уникальный ID для каждой записи
            RR_ID=$(generate_uuid)
            
            # Время с случайным сдвигом в пределах часа
            minute=$((RANDOM % 60))
            second=$((RANDOM % 60))
            TIMESTAMP=$(date_with_time $i $hour $minute $second)
            
            # Генерируем R-R интервал
            RR_INTERVAL=$(generate_rr_interval $hour)
            
            # Генерируем дополнительные метрики
            SIGNAL_QUALITY=$(generate_signal_quality $hour)
            BATTERY_LEVEL=$((80 + RANDOM % 20))  # 80-100%
            MOTION_DETECTED=$((RANDOM % 2))      # 0 или 1
            DEVICE_TEMP=$(echo "scale=1; 36.0 + $(($RANDOM % 40)) / 10" | bc -l 2>/dev/null || echo "36.5")
            
            # Подсчитываем статистику
            if [ $RR_INTERVAL -ge 300 ] && [ $RR_INTERVAL -le 2000 ]; then
                VALID_INTERVALS=$((VALID_INTERVALS + 1))
            else
                ANOMALOUS_INTERVALS=$((ANOMALOUS_INTERVALS + 1))
            fi
            
            if (( $(echo "$SIGNAL_QUALITY > 0.8" | bc -l 2>/dev/null || echo 0) )); then
                HIGH_QUALITY_COUNT=$((HIGH_QUALITY_COUNT + 1))
            fi
            
            # Записываем R-R интервал в InfluxDB
            docker exec influxdb influx write \
                --bucket $INFLUX_BUCKET \
                --org $INFLUX_ORG \
                --token $INFLUX_TOKEN \
                --precision s \
                "rr_intervals,id=$RR_ID,user_id=$USER_ID,device_id=$DEVICE_ID,quality_level=high rr_interval_ms=${RR_INTERVAL}i,signal_quality=$SIGNAL_QUALITY,battery_level=${BATTERY_LEVEL}i,motion_detected=${MOTION_DETECTED},device_temp=$DEVICE_TEMP $TIMESTAMP"
            
            TOTAL_RECORDS=$((TOTAL_RECORDS + 1))
        done
    done
done

echo ""
echo "🎉 R-R интервалы успешно добавлены!"
echo ""
echo "📈 Добавлено:"
echo "   • $TOTAL_RECORDS R-R интервалов за последние 30 дней"
echo "   • В среднем $(( TOTAL_RECORDS / 30 )) интервалов в день"
echo ""
echo "📊 Статистика качества:"
echo "   • Валидные интервалы (300-2000 мс): $VALID_INTERVALS"
echo "   • Аномальные интервалы: $ANOMALOUS_INTERVALS"
echo "   • Высокое качество сигнала (>0.8): $HIGH_QUALITY_COUNT"
echo "   • Процент валидных данных: $(( VALID_INTERVALS * 100 / TOTAL_RECORDS ))%"
echo ""
echo "🔍 Проверить данные можно:"
echo "   • В Web UI: http://localhost:8086"
echo "   • Через API приложения:"
echo "     curl -H \"Authorization: Bearer YOUR_JWT_TOKEN\" \\"
echo "          \"http://localhost:8080/api/v1/health/rr-intervals?from=2024-01-01&to=2024-12-31\""
echo ""
echo "💡 Flux запрос для проверки:"
echo "   from(bucket: \"health-data\") |> range(start: -30d) |> filter(fn: (r) => r._measurement == \"rr_intervals\")"
echo ""
echo "🏷️ Все записи содержат теги:"
echo "   • id - уникальный UUID записи"
echo "   • user_id - $USER_ID"
echo "   • device_id - $DEVICE_ID"
echo "   • quality_level - уровень качества сигнала"
echo ""
echo "📊 Поля данных:"
echo "   • rr_interval_ms - R-R интервал в миллисекундах"
echo "   • signal_quality - качество сигнала (0.0-1.0)"
echo "   • battery_level - уровень заряда (0-100)"
echo "   • motion_detected - флаг движения (0/1)"
echo "   • device_temp - температура устройства (°C)"
echo ""
echo "🖥️ Скрипт выполнен на: $OS_TYPE" 