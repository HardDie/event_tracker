-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions (
    id           INTEGER   PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER   NOT NULL REFERENCES users(id) UNIQUE,
    session_hash TEXT      NOT NULL UNIQUE,
    created_at   TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at   TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at   TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
