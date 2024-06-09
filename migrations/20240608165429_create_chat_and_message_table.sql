-- +goose Up
CREATE TABLE chat (
    id SERIAL PRIMARY KEY,
    usernames TEXT[]
);

CREATE TABLE message (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL,
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP
);

-- +goose Down
DROP TABLE chat;
DROP TABLE message;