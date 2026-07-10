-- +goose Up

-- ВНИМАНИЕ: время хранится как TIMESTAMP (без зоны).
-- Инвариант: приложение ПИШЕТ строго UTC (time.Now().UTC()), клиент конвертит в пояс юзера.
-- Любая запись не-UTC (ручной psql, миграция, вторая интеграция) молча сдвинет timeline.

CREATE TABLE chats (
    id         BIGSERIAL PRIMARY KEY,
    type       TEXT NOT NULL CHECK (type IN ('personal', 'group')),
    name       TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE TABLE chat_members (
    id      BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chats(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    UNIQUE (chat_id, user_id)
);

CREATE INDEX chat_members_user_id_idx ON chat_members (user_id);

CREATE TABLE messages (
    id         BIGSERIAL PRIMARY KEY,
    chat_id    BIGINT NOT NULL REFERENCES chats(id),
    sender_id  BIGINT NOT NULL REFERENCES users(id),
    content    TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    edited_at  TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX messages_chat_id_created_at_idx ON messages (chat_id, created_at);

-- +goose Down

DROP TABLE messages;
DROP TABLE chat_members;
DROP TABLE chats;