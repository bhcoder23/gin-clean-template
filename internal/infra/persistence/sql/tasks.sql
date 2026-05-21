-- name: CreateTask :exec
INSERT INTO tasks (
    id,
    user_id,
    title,
    description,
    status,
    created_at,
    updated_at
) VALUES (
    @id,
    @user_id,
    @title,
    @description,
    @status,
    @created_at,
    @updated_at
);

-- name: GetTaskByID :one
SELECT
    id,
    user_id,
    title,
    description,
    status,
    created_at,
    updated_at
FROM tasks
WHERE id = @id
  AND user_id = @user_id;

-- name: CountTasks :one
SELECT COUNT(*)
FROM tasks
WHERE user_id = @user_id
  AND (@status::text = '' OR status = @status::text)
  AND (@query::text = '' OR title ILIKE '%' || @query::text || '%');

-- name: ListTasks :many
SELECT
    id,
    user_id,
    title,
    description,
    status,
    created_at,
    updated_at
FROM tasks
WHERE user_id = @user_id
  AND (@status::text = '' OR status = @status::text)
  AND (@query::text = '' OR title ILIKE '%' || @query::text || '%')
ORDER BY created_at DESC
LIMIT @limit_count::bigint
OFFSET @offset_count::bigint;

-- name: UpdateTask :execrows
UPDATE tasks
SET
    title = @title,
    description = @description,
    status = @status,
    updated_at = @updated_at
WHERE id = @id
  AND user_id = @user_id;

-- name: DeleteTask :execrows
DELETE FROM tasks
WHERE id = @id
  AND user_id = @user_id;
