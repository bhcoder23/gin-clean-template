-- name: CreateNotification :exec
INSERT INTO notifications (
    id,
    user_id,
    task_id,
    type,
    title,
    body,
    read,
    created_at,
    read_at
) VALUES (
    @id,
    @user_id,
    @task_id,
    @type,
    @title,
    @body,
    @read,
    @created_at,
    @read_at
);

-- name: GetNotificationByID :one
SELECT
    id,
    user_id,
    task_id,
    type,
    title,
    body,
    read,
    created_at,
    read_at
FROM notifications
WHERE id = @id
  AND user_id = @user_id;

-- name: CountNotifications :one
SELECT COUNT(*)
FROM notifications
WHERE user_id = @user_id
  AND (NOT @unread_only::boolean OR read = FALSE);

-- name: ListNotifications :many
SELECT
    id,
    user_id,
    task_id,
    type,
    title,
    body,
    read,
    created_at,
    read_at
FROM notifications
WHERE user_id = @user_id
  AND (NOT @unread_only::boolean OR read = FALSE)
ORDER BY created_at DESC
LIMIT @limit_count::bigint
OFFSET @offset_count::bigint;

-- name: UpdateNotification :execrows
UPDATE notifications
SET
    read = @read,
    read_at = @read_at
WHERE id = @id
  AND user_id = @user_id;
