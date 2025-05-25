package health_repo

import (
	"context"
	"errors"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/repositories/dbqueries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type HealthRepo struct {
	query *dbqueries.Queries
}

func NewHealthRepo(query *dbqueries.Queries) *HealthRepo {
	return &HealthRepo{
		query: query,
	}
}

func (hr *HealthRepo) GetHeartRates(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.HeartRate, error) {
	records, err := hr.query.GetHeartRatesByUserAndDateRange(ctx, dbqueries.GetHeartRatesByUserAndDateRangeParams{
		UserID:      userID,
		CreatedAt:   from,
		CreatedAt_2: to,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	res := make([]models.HeartRate, len(records))
	for i, v := range records {
		res[i] = models.HeartRate{
			ID:        v.ID,
			UserID:    v.UserID,
			DeviceID:  v.DeviceID,
			BPM:       int64(v.Bpm),
			CreatedAt: v.CreatedAt,
		}
	}
	return res, nil
}

func (hr *HealthRepo) CreateHeartRate(ctx context.Context, heartRate models.HeartRate) error {
	_, err := hr.query.InsertHeartRate(ctx, dbqueries.InsertHeartRateParams{
		ID:        heartRate.ID,
		UserID:    heartRate.UserID,
		DeviceID:  heartRate.DeviceID,
		Bpm:       int32(heartRate.BPM),
		CreatedAt: heartRate.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (hr *HealthRepo) CreateHeartRates(ctx context.Context, heartRates []models.HeartRate) error {
	param := dbqueries.InsertHeartRatesParams{
		Column1: make([]uuid.UUID, len(heartRates)),
		Column2: make([]uuid.UUID, len(heartRates)),
		Column3: make([]uuid.UUID, len(heartRates)),
		Column4: make([]int32, len(heartRates)),
		Column5: make([]time.Time, len(heartRates)),
	}
	for i, v := range heartRates {
		param.Column1[i] = v.ID
		param.Column2[i] = v.UserID
		param.Column3[i] = v.DeviceID
		param.Column4[i] = int32(v.BPM)
		param.Column5[i] = v.CreatedAt
	}

	if err := hr.query.InsertHeartRates(ctx, param); err != nil {
		return err
	}
	return nil
}

func (hr *HealthRepo) GetSteps(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Step, error) {
	records, err := hr.query.GetStepsByUserAndDateRange(ctx, dbqueries.GetStepsByUserAndDateRangeParams{
		UserID:      userID,
		CreatedAt:   from,
		CreatedAt_2: to,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	res := make([]models.Step, len(records))
	for i, v := range records {
		res[i] = models.Step{
			ID:        v.ID,
			UserID:    v.UserID,
			DeviceID:  v.DeviceID,
			StepCount: int64(v.StepCount),
			CreatedAt: v.CreatedAt,
		}
	}
	return res, nil
}

func (hr *HealthRepo) CreateStep(ctx context.Context, step models.Step) error {
	_, err := hr.query.InsertStep(ctx, dbqueries.InsertStepParams{
		ID:        step.ID,
		UserID:    step.UserID,
		DeviceID:  step.DeviceID,
		StepCount: int32(step.StepCount),
		CreatedAt: step.CreatedAt,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == "STEPS_PK":
				return apperrors.ErrAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (hr *HealthRepo) CreateSteps(ctx context.Context, steps []models.Step) error {
	param := dbqueries.InsertStepsParams{
		Column1: make([]uuid.UUID, len(steps)),
		Column2: make([]uuid.UUID, len(steps)),
		Column3: make([]uuid.UUID, len(steps)),
		Column4: make([]int32, len(steps)),
		Column5: make([]time.Time, len(steps)),
	}
	for i, v := range steps {
		param.Column1[i] = v.ID
		param.Column2[i] = v.UserID
		param.Column3[i] = v.DeviceID
		param.Column4[i] = int32(v.StepCount)
		param.Column5[i] = v.CreatedAt
	}

	if err := hr.query.InsertSteps(ctx, param); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == "STEPS_PK":
				return apperrors.ErrAlreadyExists
			case pgErr.ConstraintName == "UNIQUE_STEPS":
				return apperrors.ErrAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (hr *HealthRepo) GetSleeps(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Sleep, error) {
	records, err := hr.query.GetSleepsByUserAndDateRange(ctx, dbqueries.GetSleepsByUserAndDateRangeParams{
		UserID:      userID,
		StartedAt:   from,
		StartedAt_2: to,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	res := make([]models.Sleep, len(records))
	for i, v := range records {
		res[i] = models.Sleep{
			ID:        v.ID,
			UserID:    v.UserID,
			DeviceID:  v.DeviceID,
			StartedAt: v.StartedAt,
			EndedAt:   v.EndedAt,
		}
	}
	return res, nil
}

func (hr *HealthRepo) CreateSleep(ctx context.Context, sleep models.Sleep) error {
	_, err := hr.query.InsertSleep(ctx, dbqueries.InsertSleepParams{
		ID:        sleep.ID,
		UserID:    sleep.UserID,
		DeviceID:  sleep.DeviceID,
		StartedAt: sleep.StartedAt,
		EndedAt:   sleep.EndedAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (hr *HealthRepo) CreateSleeps(ctx context.Context, sleeps []models.Sleep) error {
	param := dbqueries.InsertSleepsParams{
		Column1: make([]uuid.UUID, len(sleeps)),
		Column2: make([]uuid.UUID, len(sleeps)),
		Column3: make([]uuid.UUID, len(sleeps)),
		Column4: make([]time.Time, len(sleeps)),
		Column5: make([]time.Time, len(sleeps)),
	}
	for i, v := range sleeps {
		param.Column1[i] = v.ID
		param.Column2[i] = v.UserID
		param.Column3[i] = v.DeviceID
		param.Column4[i] = v.StartedAt
		param.Column5[i] = v.EndedAt
	}

	if err := hr.query.InsertSleeps(ctx, param); err != nil {
		return err
	}
	return nil
}

func (hr *HealthRepo) GetTemperatures(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Temperature, error) {
	records, err := hr.query.GetTemperaturesByUserAndDateRange(ctx, dbqueries.GetTemperaturesByUserAndDateRangeParams{
		UserID:      userID,
		CreatedAt:   from,
		CreatedAt_2: to,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	res := make([]models.Temperature, len(records))
	for i, v := range records {
		res[i] = models.Temperature{
			ID:                 v.ID,
			UserID:             v.UserID,
			DeviceID:           v.DeviceID,
			TemperatureCelsius: v.TemperatureCelsius,
			CreatedAt:          v.CreatedAt,
		}
	}
	return res, nil
}

func (hr *HealthRepo) CreateTemperature(ctx context.Context, temperature models.Temperature) error {
	_, err := hr.query.InsertTemperature(ctx, dbqueries.InsertTemperatureParams{
		ID:                 temperature.ID,
		UserID:             temperature.UserID,
		DeviceID:           temperature.DeviceID,
		TemperatureCelsius: temperature.TemperatureCelsius,
		CreatedAt:          temperature.CreatedAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (hr *HealthRepo) CreateTemperatures(ctx context.Context, temperatures []models.Temperature) error {
	param := dbqueries.InsertTemperaturesParams{
		Column1: make([]uuid.UUID, len(temperatures)),
		Column2: make([]uuid.UUID, len(temperatures)),
		Column3: make([]uuid.UUID, len(temperatures)),
		Column4: make([]float64, len(temperatures)),
		Column5: make([]time.Time, len(temperatures)),
	}
	for i, v := range temperatures {
		param.Column1[i] = v.ID
		param.Column2[i] = v.UserID
		param.Column3[i] = v.DeviceID
		param.Column4[i] = v.TemperatureCelsius
		param.Column5[i] = v.CreatedAt
	}

	if err := hr.query.InsertTemperatures(ctx, param); err != nil {
		return err
	}

	return nil
}

func (hr *HealthRepo) GetWeights(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]models.Weight, error) {
	records, err := hr.query.GetWeightsByUserAndDateRange(ctx, dbqueries.GetWeightsByUserAndDateRangeParams{
		UserID:      userID,
		CreatedAt:   from,
		CreatedAt_2: to,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	res := make([]models.Weight, len(records))
	for i, v := range records {
		res[i] = models.Weight{
			ID:        v.ID,
			UserID:    v.UserID,
			DeviceID:  v.DeviceID,
			WeightKg:  v.WeightKg,
			CreatedAt: v.CreatedAt,
		}
	}
	return res, nil
}

func (hr *HealthRepo) CreateWeight(ctx context.Context, weight models.Weight) error {
	_, err := hr.query.InsertWeight(ctx, dbqueries.InsertWeightParams{
		ID:        weight.ID,
		UserID:    weight.UserID,
		DeviceID:  weight.DeviceID,
		WeightKg:  weight.WeightKg,
		CreatedAt: weight.CreatedAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (hr *HealthRepo) CreateWeights(ctx context.Context, weights []models.Weight) error {
	param := dbqueries.InsertWeightsParams{
		Column1: make([]uuid.UUID, len(weights)),
		Column2: make([]uuid.UUID, len(weights)),
		Column3: make([]uuid.UUID, len(weights)),
		Column4: make([]float64, len(weights)),
		Column5: make([]time.Time, len(weights)),
	}
	for i, v := range weights {
		param.Column1[i] = v.ID
		param.Column2[i] = v.UserID
		param.Column3[i] = v.DeviceID
		param.Column4[i] = v.WeightKg
		param.Column5[i] = v.CreatedAt
	}

	if err := hr.query.InsertWeights(ctx, param); err != nil {
		return err
	}
	return nil
}
