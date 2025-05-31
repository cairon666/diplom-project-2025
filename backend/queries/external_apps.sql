-- name: CreateExternalApp :one
INSERT INTO "EXTERNAL_APPS" (id, name, owner_user_id, api_key_hash, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetExternalAppByID :one
SELECT * FROM "EXTERNAL_APPS"
WHERE id = $1;

-- name: GetExternalAppByAPIKeyHash :one
SELECT * FROM "EXTERNAL_APPS"
WHERE api_key_hash = $1;

-- name: ListExternalAppsByOwner :many
SELECT * FROM "EXTERNAL_APPS"
WHERE owner_user_id = $1
ORDER BY created_at DESC;

-- name: DeleteExternalApp :one
DELETE FROM "EXTERNAL_APPS"
WHERE id = $1
RETURNING id;

-- name: UpdateExternalAppName :exec
UPDATE "EXTERNAL_APPS"
SET name = $2
WHERE id = $1;
