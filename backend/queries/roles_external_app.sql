-- name: AddRoleToExternalApp :exec
INSERT INTO "EXTERNAL_APPS_ROLES" (external_app_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveRoleFromExternalApp :exec
DELETE FROM "EXTERNAL_APPS_ROLES"
WHERE external_app_id = $1 AND role_id = $2;

-- name: GetRolesByExternalAppID :many
SELECT r.*
FROM "ROLES" r
         JOIN "EXTERNAL_APPS_ROLES" ear ON r.id = ear.role_id
WHERE ear.external_app_id = $1;

-- name: GetPermissionsByExternalAppID :many
SELECT p.*
FROM "PERMISSIONS" p
         JOIN "ROLE_PERMISSIONS" rp ON p.id = rp.permission_id
         JOIN "EXTERNAL_APPS_ROLES" ear ON rp.role_id = ear.role_id
WHERE ear.external_app_id = $1;

-- name: ExternalAppHasRole :one
SELECT EXISTS (
    SELECT 1 FROM "EXTERNAL_APPS_ROLES"
    WHERE external_app_id = $1 AND role_id = $2
) AS exists;
