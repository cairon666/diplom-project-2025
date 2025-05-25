-- name: HasUserPermission :one
SELECT EXISTS (
    SELECT 1
    FROM "USER_ROLES" ur
             JOIN "ROLE_PERMISSIONS" rp ON ur.role_id = rp.role_id
             JOIN "PERMISSIONS" p ON rp.permission_id = p.id
    WHERE ur.user_id = $1
      AND p.name = $2
) AS has_permission;

-- name: AssignRoleToUser :exec
INSERT INTO "USER_ROLES" (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveRoleFromUser :exec
DELETE FROM "USER_ROLES" WHERE user_id = $1 AND role_id = $2;

-- name: GetPermissionsByUserID :many
SELECT DISTINCT p.id, p.name
FROM "USER_ROLES" ur
         JOIN "ROLE_PERMISSIONS" rp ON ur.role_id = rp.role_id
         JOIN "PERMISSIONS" p ON rp.permission_id = p.id
WHERE ur.user_id = $1;