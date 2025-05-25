-- name: GetUserById :one
SELECT *
from "USERS"
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
from "USERS"
WHERE email = $1;

-- name: CreateUser :exec
INSERT INTO "USERS" (id, email, first_name, last_name, is_registration_complete, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUserEmail :exec
UPDATE "USERS"
SET email = $2
WHERE id = $1;

-- name: UpdateUserData :exec
UPDATE "USERS"
SET (first_name, last_name) = ($2, $3)
WHERE id = $1;

-- name: UpdateUserFull :exec
UPDATE "USERS"
SET (first_name, last_name, email, is_registration_complete) = ($2, $3, $4, $5)
WHERE id = $1;