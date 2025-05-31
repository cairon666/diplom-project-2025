-- name: GetUserPassword :one
SELECT *
from "USER_PASSWORDS"
WHERE user_id = $1;

-- name: CreateUserPassword :exec
INSERT INTO "USER_PASSWORDS" (user_id, password_hash, salt)
VALUES ($1, $2, $3);

-- name: UpdateUserPassword :exec
UPDATE "USER_PASSWORDS"
SET password_hash = $2, salt = $3
WHERE user_id = $1;

-- name: DeleteUserPassword :exec
DELETE
FROM "USER_PASSWORDS"
WHERE user_id = $1;