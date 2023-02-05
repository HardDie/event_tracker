-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS passwords (
    id              SERIAL    PRIMARY KEY,
    user_id         INT       NOT NULL UNIQUE REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    password_hash   TEXT      NOT NULL,
    failed_attempts INT       NOT NULL DEFAULT (0),
    created_at      TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at      TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at      TIMESTAMP,
    blocked_at      TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE passwords;
-- +goose StatementEnd
