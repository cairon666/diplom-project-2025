-- STEP QUERIES
-- name: InsertStep :one
INSERT INTO "STEPS" (id, user_id, device_id, step_count, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: InsertSteps :exec
INSERT INTO "STEPS" (id, user_id, device_id, step_count, created_at)
SELECT UNNEST($1::uuid[]), UNNEST($2::uuid[]), UNNEST($3::uuid[]), UNNEST($4::integer[]), UNNEST($5::timestamptz[]);

-- name: GetStepsByUserAndDateRange :many
SELECT * FROM "STEPS"
WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
ORDER BY created_at DESC;

-- HEART RATE QUERIES
-- name: InsertHeartRate :one
INSERT INTO "HEART_RATES" (id, user_id, device_id, bpm, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: InsertHeartRates :exec
INSERT INTO "HEART_RATES" (id, user_id, device_id, bpm, created_at)
SELECT UNNEST($1::uuid[]), UNNEST($2::uuid[]), UNNEST($3::uuid[]), UNNEST($4::integer[]), UNNEST($5::timestamptz[]);

-- name: GetHeartRatesByUserAndDateRange :many
SELECT * FROM "HEART_RATES"
WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
ORDER BY created_at DESC;

-- TEMPERATURE QUERIES
-- name: InsertTemperature :one
INSERT INTO "TEMPERATURES" (id, user_id, device_id, temperature_celsius, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: InsertTemperatures :exec
INSERT INTO "TEMPERATURES" (id, user_id, device_id, temperature_celsius, created_at)
SELECT UNNEST($1::uuid[]), UNNEST($2::uuid[]), UNNEST($3::uuid[]), UNNEST($4::float[]), UNNEST($5::timestamptz[]);

-- name: GetTemperaturesByUserAndDateRange :many
SELECT * FROM "TEMPERATURES"
WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
ORDER BY created_at DESC;

-- WEIGHT QUERIES
-- name: InsertWeight :one
INSERT INTO "WEIGHTS" (id, user_id, device_id, weight_kg, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: InsertWeights :exec
INSERT INTO "WEIGHTS" (id, user_id, device_id, weight_kg, created_at)
SELECT UNNEST($1::uuid[]), UNNEST($2::uuid[]), UNNEST($3::uuid[]), UNNEST($4::float[]), UNNEST($5::timestamptz[]);

-- name: GetWeightsByUserAndDateRange :many
SELECT * FROM "WEIGHTS"
WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
ORDER BY created_at DESC;

-- SLEEP QUERIES
-- name: InsertSleep :one
INSERT INTO "SLEEPS" (id, user_id, device_id, started_at, ended_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: InsertSleeps :exec
INSERT INTO "SLEEPS" (id, user_id, device_id, started_at, ended_at)
SELECT UNNEST($1::uuid[]), UNNEST($2::uuid[]), UNNEST($3::uuid[]), UNNEST($4::timestamptz[]), UNNEST($5::timestamptz[]);

-- name: GetSleepsByUserAndDateRange :many
SELECT * FROM "SLEEPS"
WHERE user_id = $1 AND started_at BETWEEN $2 AND $3
ORDER BY started_at DESC;
