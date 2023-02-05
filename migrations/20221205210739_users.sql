-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id              SERIAL    PRIMARY KEY,
    username        TEXT      NOT NULL UNIQUE,
    displayed_name  TEXT,
    email           TEXT,
    profile_image   TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at      TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at      TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
