-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS friends (
    id           INTEGER   PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER   NOT NULL REFERENCES users(id),
    with_user_id INTEGER   NOT NULL REFERENCES users(id),
    created_at   TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at   TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at   TIMESTAMP
);
CREATE INDEX friends_id_idx ON friends (id);
CREATE INDEX friends_user_id_idx ON friends (user_id);
CREATE INDEX friends_with_user_id_idx ON friends (with_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friends;
-- +goose StatementEnd
