-- name: CreateChat :one
INSERT INTO chats (type, name)
VALUES ($1, $2)
RETURNING id, type, name, created_at;

-- name: GetChatByID :one
SELECT id, type, name, created_at
FROM chats
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindPersonalChat :one
SELECT c.id
FROM chats c
JOIN chat_members m1 ON m1.chat_id = c.id AND m1.user_id = $1
JOIN chat_members m2 on m2.chat_id = c.id AND m2.user_id = $2
WHERE c.type = 'personal' AND c.deleted_at IS NULL
LIMIT 1;

-- name: SoftDeleteChat :exec
UPDATE chats SET deleted_at = now() WHERE id = $1;

-- name: ListUserChats :many
SELECT
    c.id,
    c.type,
    c.name,
    other.id        AS other_user_id,
    other.username  AS other_username
FROM chats c
JOIN chat_members me ON me.chat_id = c.id AND me.user_id = $1
LEFT JOIN LATERAL (
    SELECT u.id, u.username
    FROM chat_members cm
    JOIN users u ON u.id = cm.user_id
    WHERE cm.chat_id = c.id
        AND cm.user_id <> $1
        AND c.type = 'personal'
    LIMIT 1
) other ON true
WHERE c.deleted_at IS NULL
ORDER BY c.id DESC;