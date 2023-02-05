-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions (
    id           SERIAL    PRIMARY KEY,
    user_id      INT       NOT NULL UNIQUE REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    session_hash TEXT      NOT NULL UNIQUE,
    created_at   TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at   TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at   TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
