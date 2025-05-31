package health_repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/config"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/pkg/influxdb"
	"github.com/google/uuid"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// HealthRepo реализует репозиторий для работы с данными здоровья в InfluxDB
type HealthRepo struct {
	influxClient influxdb.InfluxDBClient
	org          string
	bucket       string
}

// NewHealthRepo создает новый экземпляр репозитория
func NewHealthRepo(influxClient influxdb.InfluxDBClient, conf *config.Config) *HealthRepo {
	return &HealthRepo{
		influxClient: influxClient,
		org:          conf.InfluxDB.Org,
		bucket:       conf.InfluxDB.Bucket,
	}
}

// GetSteps получает данные о шагах пользователя за период
func (r *HealthRepo) GetSteps(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Step, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "steps")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "count")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.HealthDataQueryFailedf("failed to query steps: %v", err)
	}
	defer result.Close()

	var steps []models.Step
	for result.Next() {
		record := result.Record()

		// Извлекаем ID из tags
		idStr, ok := record.ValueByKey("id").(string)
		if !ok {
			continue
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		deviceIDStr, ok := record.ValueByKey("device_id").(string)
		if !ok {
			continue
		}
		deviceID, err := uuid.Parse(deviceIDStr)
		if err != nil {
			continue
		}

		count, ok := record.Value().(int64)
		if !ok {
			continue
		}

		step := models.Step{
			ID:        id,
			UserID:    userID,
			DeviceID:  deviceID,
			StepCount: count,
			CreatedAt: record.Time(),
		}
		steps = append(steps, step)
	}

	if result.Err() != nil {
		return nil, apperrors.HealthDataReadFailedf("error reading steps result: %v", result.Err())
	}

	return steps, nil
}

// GetHeartRates получает данные о пульсе пользователя за период
func (r *HealthRepo) GetHeartRates(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.HeartRate, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "heart_rate")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "bpm")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.HealthDataQueryFailedf("failed to query heart rates: %v", err)
	}
	defer result.Close()

	var heartRates []models.HeartRate
	for result.Next() {
		record := result.Record()

		// Извлекаем ID из tags
		idStr, ok := record.ValueByKey("id").(string)
		if !ok {
			continue
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		deviceIDStr, ok := record.ValueByKey("device_id").(string)
		if !ok {
			continue
		}
		deviceID, err := uuid.Parse(deviceIDStr)
		if err != nil {
			continue
		}

		bpm, ok := record.Value().(int64)
		if !ok {
			continue
		}

		heartRate := models.HeartRate{
			ID:        id,
			UserID:    userID,
			DeviceID:  deviceID,
			BPM:       bpm,
			CreatedAt: record.Time(),
		}
		heartRates = append(heartRates, heartRate)
	}

	if result.Err() != nil {
		return nil, apperrors.HealthDataReadFailedf("error reading heart rates result: %v", result.Err())
	}

	return heartRates, nil
}

// GetWeights получает данные о весе пользователя за период
func (r *HealthRepo) GetWeights(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Weight, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "weight")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "weight_kg")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.HealthDataQueryFailedf("failed to query weights: %v", err)
	}
	defer result.Close()

	var weights []models.Weight
	for result.Next() {
		record := result.Record()

		// Извлекаем ID из tags
		idStr, ok := record.ValueByKey("id").(string)
		if !ok {
			continue
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		deviceIDStr, ok := record.ValueByKey("device_id").(string)
		if !ok {
			continue
		}
		deviceID, err := uuid.Parse(deviceIDStr)
		if err != nil {
			continue
		}

		weightKg, ok := record.Value().(float64)
		if !ok {
			continue
		}

		weight := models.Weight{
			ID:        id,
			UserID:    userID,
			DeviceID:  deviceID,
			WeightKg:  weightKg,
			CreatedAt: record.Time(),
		}
		weights = append(weights, weight)
	}

	if result.Err() != nil {
		return nil, apperrors.HealthDataReadFailedf("error reading weights result: %v", result.Err())
	}

	return weights, nil
}

// GetTemperatures получает данные о температуре пользователя за период
func (r *HealthRepo) GetTemperatures(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Temperature, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "temperature")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "temperature_celsius")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.HealthDataQueryFailedf("failed to query temperatures: %v", err)
	}
	defer result.Close()

	var temperatures []models.Temperature
	for result.Next() {
		record := result.Record()

		// Извлекаем ID из tags
		idStr, ok := record.ValueByKey("id").(string)
		if !ok {
			continue
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		deviceIDStr, ok := record.ValueByKey("device_id").(string)
		if !ok {
			continue
		}
		deviceID, err := uuid.Parse(deviceIDStr)
		if err != nil {
			continue
		}

		tempCelsius, ok := record.Value().(float64)
		if !ok {
			continue
		}

		temperature := models.Temperature{
			ID:                 id,
			UserID:             userID,
			DeviceID:           deviceID,
			TemperatureCelsius: tempCelsius,
			CreatedAt:          record.Time(),
		}
		temperatures = append(temperatures, temperature)
	}

	if result.Err() != nil {
		return nil, apperrors.HealthDataReadFailedf("error reading temperatures result: %v", result.Err())
	}

	return temperatures, nil
}

// GetSleeps получает данные о сне пользователя за период
func (r *HealthRepo) GetSleeps(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Sleep, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "sleep")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "start_time" or r._field == "end_time")
		|> pivot(rowKey: ["_time", "id"], columnKey: ["_field"], valueColumn: "_value")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, apperrors.HealthDataQueryFailedf("failed to query sleeps: %v", err)
	}
	defer result.Close()

	var sleeps []models.Sleep
	for result.Next() {
		record := result.Record()

		// Извлекаем ID из tags
		idStr, ok := record.ValueByKey("id").(string)
		if !ok {
			continue
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		deviceIDStr, ok := record.ValueByKey("device_id").(string)
		if !ok {
			continue
		}
		deviceID, err := uuid.Parse(deviceIDStr)
		if err != nil {
			continue
		}

		// Получаем start_time из pivot результата
		startTimeValue := record.ValueByKey("start_time")
		if startTimeValue == nil {
			continue
		}
		startTimeUnix, ok := startTimeValue.(int64)
		if !ok {
			// Попробуем преобразовать из float64
			if startTimeFloat, ok := startTimeValue.(float64); ok {
				startTimeUnix = int64(startTimeFloat)
			} else {
				continue
			}
		}
		startTime := time.Unix(startTimeUnix, 0)

		// Получаем end_time из pivot результата
		endTimeValue := record.ValueByKey("end_time")
		if endTimeValue == nil {
			continue
		}
		endTimeUnix, ok := endTimeValue.(int64)
		if !ok {
			// Попробуем преобразовать из float64
			if endTimeFloat, ok := endTimeValue.(float64); ok {
				endTimeUnix = int64(endTimeFloat)
			} else {
				continue
			}
		}
		endTime := time.Unix(endTimeUnix, 0)

		sleep := models.Sleep{
			ID:        id,
			UserID:    userID,
			DeviceID:  deviceID,
			StartedAt: startTime,
			EndedAt:   endTime,
		}
		sleeps = append(sleeps, sleep)
	}

	if result.Err() != nil {
		return nil, apperrors.HealthDataReadFailedf("error reading sleeps result: %v", result.Err())
	}

	return sleeps, nil
}

// Методы для записи данных

// CreateStep создает запись о шагах
func (r *HealthRepo) CreateStep(ctx context.Context, step models.Step) error {
	writeAPI := r.influxClient.WriteAPIBlocking(r.org, r.bucket)

	point := write.NewPoint("steps",
		map[string]string{
			"id":        step.ID.String(), // ID для идемпотентности
			"user_id":   step.UserID.String(),
			"device_id": step.DeviceID.String(),
		},
		map[string]interface{}{
			"count": step.StepCount,
		},
		step.CreatedAt)

	return writeAPI.WritePoint(ctx, point)
}

// CreateHeartRate создает запись о пульсе
func (r *HealthRepo) CreateHeartRate(ctx context.Context, heartRate models.HeartRate) error {
	writeAPI := r.influxClient.WriteAPIBlocking(r.org, r.bucket)

	point := write.NewPoint("heart_rate",
		map[string]string{
			"id":        heartRate.ID.String(), // ID для идемпотентности
			"user_id":   heartRate.UserID.String(),
			"device_id": heartRate.DeviceID.String(),
		},
		map[string]interface{}{
			"bpm": heartRate.BPM,
		},
		heartRate.CreatedAt)

	return writeAPI.WritePoint(ctx, point)
}

// CreateWeight создает запись о весе
func (r *HealthRepo) CreateWeight(ctx context.Context, weight models.Weight) error {
	writeAPI := r.influxClient.WriteAPIBlocking(r.org, r.bucket)

	point := write.NewPoint("weight",
		map[string]string{
			"id":        weight.ID.String(), // ID для идемпотентности
			"user_id":   weight.UserID.String(),
			"device_id": weight.DeviceID.String(),
		},
		map[string]interface{}{
			"weight_kg": weight.WeightKg,
		},
		weight.CreatedAt)

	return writeAPI.WritePoint(ctx, point)
}

// CreateTemperature создает запись о температуре
func (r *HealthRepo) CreateTemperature(ctx context.Context, temperature models.Temperature) error {
	writeAPI := r.influxClient.WriteAPIBlocking(r.org, r.bucket)

	point := write.NewPoint("temperature",
		map[string]string{
			"id":        temperature.ID.String(), // ID для идемпотентности
			"user_id":   temperature.UserID.String(),
			"device_id": temperature.DeviceID.String(),
		},
		map[string]interface{}{
			"temperature_celsius": temperature.TemperatureCelsius,
		},
		temperature.CreatedAt)

	return writeAPI.WritePoint(ctx, point)
}

// CreateSleep создает запись о сне
func (r *HealthRepo) CreateSleep(ctx context.Context, sleep models.Sleep) error {
	writeAPI := r.influxClient.WriteAPIBlocking(r.org, r.bucket)

	point := write.NewPoint("sleep",
		map[string]string{
			"id":        sleep.ID.String(), // ID для идемпотентности
			"user_id":   sleep.UserID.String(),
			"device_id": sleep.DeviceID.String(),
		},
		map[string]interface{}{
			"start_time": sleep.StartedAt.Unix(),
			"end_time":   sleep.EndedAt.Unix(),
		},
		sleep.StartedAt)

	return writeAPI.WritePoint(ctx, point)
}

// Batch методы для массовой записи

// CreateSteps создает множественные записи о шагах
func (r *HealthRepo) CreateSteps(ctx context.Context, steps []models.Step) error {
	for _, step := range steps {
		if err := r.CreateStep(ctx, step); err != nil {
			return err
		}
	}

	return nil
}

// CreateHeartRates создает множественные записи о пульсе
func (r *HealthRepo) CreateHeartRates(ctx context.Context, heartRates []models.HeartRate) error {
	for _, hr := range heartRates {
		if err := r.CreateHeartRate(ctx, hr); err != nil {
			return err
		}
	}

	return nil
}

// CreateWeights создает множественные записи о весе
func (r *HealthRepo) CreateWeights(ctx context.Context, weights []models.Weight) error {
	for _, weight := range weights {
		if err := r.CreateWeight(ctx, weight); err != nil {
			return err
		}
	}

	return nil
}

// CreateTemperatures создает множественные записи о температуре
func (r *HealthRepo) CreateTemperatures(ctx context.Context, temperatures []models.Temperature) error {
	for _, temp := range temperatures {
		if err := r.CreateTemperature(ctx, temp); err != nil {
			return err
		}
	}

	return nil
}

// CreateSleeps создает множественные записи о сне
func (r *HealthRepo) CreateSleeps(ctx context.Context, sleeps []models.Sleep) error {
	for _, sleep := range sleeps {
		if err := r.CreateSleep(ctx, sleep); err != nil {
			return err
		}
	}

	return nil
}

// Методы агрегации данных

// GetHourlySteps возвращает агрегированные данные по шагам по часам
func (r *HealthRepo) GetHourlySteps(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]int64, error) {
	// Валидация времени
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	// Если диапазон слишком маленький, возвращаем пустой результат
	if to.Sub(from) < time.Minute {
		return make(map[time.Time]int64), nil
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "steps")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "count")
		|> aggregateWindow(every: 1h, fn: sum, createEmpty: false)
		|> yield(name: "hourly_steps")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	return r.executeAggregationQuery(ctx, query)
}

// GetDailySteps возвращает агрегированные данные по шагам по дням
func (r *HealthRepo) GetDailySteps(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]int64, error) {
	// Валидация времени
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	// Если диапазон слишком маленький, возвращаем пустой результат
	if to.Sub(from) < time.Minute {
		return make(map[time.Time]int64), nil
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "steps")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "count")
		|> aggregateWindow(every: 1d, fn: sum, createEmpty: false)
		|> yield(name: "daily_steps")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	return r.executeAggregationQuery(ctx, query)
}

// GetHourlyHeartRateAvg возвращает средний пульс по часам
func (r *HealthRepo) GetHourlyHeartRateAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "heart_rate")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "bpm")
		|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
		|> yield(name: "hourly_heart_rate")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	return r.executeFloatAggregationQuery(ctx, query)
}

// GetDailyHeartRateAvg возвращает средний пульс по дням
func (r *HealthRepo) GetDailyHeartRateAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	// Валидация времени
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	// Если диапазон слишком маленький, возвращаем пустой результат
	if to.Sub(from) < time.Minute {
		return make(map[time.Time]float64), nil
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "heart_rate")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "bpm")
		|> aggregateWindow(every: 1d, fn: mean, createEmpty: false)
		|> yield(name: "daily_heart_rate")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	return r.executeFloatAggregationQuery(ctx, query)
}

// GetDailyWeightAvg возвращает средний вес по дням
func (r *HealthRepo) GetDailyWeightAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	// Валидация времени
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	// Если диапазон слишком маленький, возвращаем пустой результат
	if to.Sub(from) < time.Minute {
		return make(map[time.Time]float64), nil
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "weight")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "weight_kg")
		|> aggregateWindow(every: 1d, fn: mean, createEmpty: false)
		|> yield(name: "daily_weight")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	return r.executeFloatAggregationQuery(ctx, query)
}

// GetHourlyTemperatureAvg возвращает среднюю температуру по часам
func (r *HealthRepo) GetHourlyTemperatureAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "temperature")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "temperature_celsius")
		|> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
		|> yield(name: "hourly_temperature")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	return r.executeFloatAggregationQuery(ctx, query)
}

// GetDailyTemperatureAvg возвращает среднюю температуру по дням
func (r *HealthRepo) GetDailyTemperatureAvg(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "temperature")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "temperature_celsius")
		|> aggregateWindow(every: 1d, fn: mean, createEmpty: false)
		|> yield(name: "daily_temperature")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	return r.executeFloatAggregationQuery(ctx, query)
}

// GetDailySleepDuration возвращает продолжительность сна по дням
func (r *HealthRepo) GetDailySleepDuration(ctx context.Context, userID uuid.UUID, from, to time.Time) (map[time.Time]float64, error) {
	// Валидация времени
	if from.After(to) {
		return nil, apperrors.InvalidTimeRangef("from time (%v) cannot be after to time (%v)", from, to)
	}

	// Если диапазон слишком маленький, возвращаем пустой результат
	if to.Sub(from) < time.Minute {
		return make(map[time.Time]float64), nil
	}

	query := fmt.Sprintf(`
		import "date"
		
		sleep_data = from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "sleep")
		|> filter(fn: (r) => r.user_id == "%s")
		|> filter(fn: (r) => r._field == "start_time" or r._field == "end_time")
		|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
		|> map(fn: (r) => ({
			_time: r._time,
			_value: float(v: r.end_time - r.start_time) / 3600.0
		}))
		|> aggregateWindow(every: 1d, fn: sum, createEmpty: false)
		|> yield(name: "daily_sleep_duration")
	`, r.bucket, from.Format(time.RFC3339), to.Format(time.RFC3339), userID.String())

	return r.executeFloatAggregationQuery(ctx, query)
}

// executeAggregationQuery выполняет запрос агрегации и возвращает результат как map[time.Time]int64
func (r *HealthRepo) executeAggregationQuery(ctx context.Context, query string) (map[time.Time]int64, error) {
	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		// Проверяем специфичные ошибки InfluxDB
		if strings.Contains(err.Error(), "cannot query an empty range") {
			// Возвращаем пустой результат вместо ошибки
			return make(map[time.Time]int64), nil
		}
		return nil, apperrors.HealthDataQueryFailedf("failed to execute aggregation query: %v", err)
	}
	defer result.Close()

	data := make(map[time.Time]int64)
	for result.Next() {
		record := result.Record()

		value, ok := record.Value().(int64)
		if !ok {
			// Попробуем преобразовать float64 в int64
			if floatVal, ok := record.Value().(float64); ok {
				value = int64(floatVal)
			} else {
				continue
			}
		}

		data[record.Time()] = value
	}

	if result.Err() != nil {
		// Проверяем специфичные ошибки InfluxDB
		if strings.Contains(result.Err().Error(), "cannot query an empty range") {
			return make(map[time.Time]int64), nil
		}
		return nil, apperrors.HealthDataReadFailedf("error reading aggregation result: %v", result.Err())
	}

	return data, nil
}

// executeFloatAggregationQuery выполняет запрос агрегации и возвращает результат как map[time.Time]float64
func (r *HealthRepo) executeFloatAggregationQuery(ctx context.Context, query string) (map[time.Time]float64, error) {
	queryAPI := r.influxClient.QueryAPI(r.org)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		// Проверяем специфичные ошибки InfluxDB
		if strings.Contains(err.Error(), "cannot query an empty range") {
			// Возвращаем пустой результат вместо ошибки
			return make(map[time.Time]float64), nil
		}
		return nil, apperrors.HealthDataQueryFailedf("failed to execute float aggregation query: %v", err)
	}
	defer result.Close()

	data := make(map[time.Time]float64)
	for result.Next() {
		record := result.Record()

		value, ok := record.Value().(float64)
		if !ok {
			// Попробуем преобразовать int64 в float64
			if intVal, ok := record.Value().(int64); ok {
				value = float64(intVal)
			} else {
				continue
			}
		}

		data[record.Time()] = value
	}

	if result.Err() != nil {
		// Проверяем специфичные ошибки InfluxDB
		if strings.Contains(result.Err().Error(), "cannot query an empty range") {
			return make(map[time.Time]float64), nil
		}
		return nil, apperrors.HealthDataReadFailedf("error reading float aggregation result: %v", result.Err())
	}

	return data, nil
}