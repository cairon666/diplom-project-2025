-- name: GetAuthProviderById :one
SELECT *
from "AUTH_PROVIDERS"
WHERE id = $1;

-- name: GetAuthProvidersByUserId :many
SELECT *
from "AUTH_PROVIDERS"
WHERE user_id = $1;

-- name: GetAuthProviderByProviderUserIdAndProviderName :one
SELECT *
from "AUTH_PROVIDERS"
WHERE provider_user_id = $1 AND provider_name = $2;

-- name: GetAuthProviderByUserIdAndProviderName :one
SELECT *
from "AUTH_PROVIDERS"
WHERE user_id = $1 AND provider_name = $2;

-- name: CreateAuthProvider :exec
INSERT INTO "AUTH_PROVIDERS" (id, user_id, provider_name, provider_user_id, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: DeleteAuthProviderById :exec
DELETE
FROM "AUTH_PROVIDERS"
WHERE id = $1;

-- name: DeleteAuthProviderByUserId :exec
DELETE
FROM "AUTH_PROVIDERS"
WHERE user_id = $1;