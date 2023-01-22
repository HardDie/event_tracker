-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS friend_invites (
    id           INTEGER   PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER   NOT NULL REFERENCES users(id),
    with_user_id INTEGER   NOT NULL REFERENCES users(id),
    created_at   TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at   TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at   TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friend_invites;
-- +goose StatementEnd
