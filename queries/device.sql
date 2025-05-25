-- name: CreateDevice :one
INSERT INTO "DEVICES" (id, user_id, device_name, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDeviceByID :one
SELECT *
FROM "DEVICES"
WHERE id = $1;

-- name: ListDevicesByUserID :many
SELECT *
FROM "DEVICES"
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateDeviceName :one
UPDATE "DEVICES"
SET device_name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteDevice :one
DELETE FROM "DEVICES"
WHERE id = $1
RETURNING id;
