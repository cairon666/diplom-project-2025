-- name: GetRolesByUserID :many
SELECT r.id, r.name
FROM "USER_ROLES" ur
         JOIN "ROLES" r ON ur.role_id = r.id
WHERE ur.user_id = $1;

-- name: GetPermissionsByRoleIDs :many
SELECT p.id, p.name
FROM "ROLE_PERMISSIONS" rp
         JOIN "PERMISSIONS" p ON rp.permission_id = p.id
WHERE rp.role_id = ANY($1::int[]);



-- name: GetRoleByName :one
SELECT id, name FROM "ROLES" WHERE name = $1;