-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS passwords (
    id              INTEGER   PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER   NOT NULL REFERENCES users(id) UNIQUE,
    password_hash   TEXT      NOT NULL,
    failed_attempts INTEGER   NOT NULL DEFAULT (0),
    created_at      TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at      TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at      TIMESTAMP,
    blocked_at      TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE passwords;
-- +goose StatementEnd
