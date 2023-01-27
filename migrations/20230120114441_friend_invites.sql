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
CREATE INDEX friend_invites_id_idx ON friend_invites (id);
CREATE INDEX friend_invites_user_id_idx ON friend_invites (user_id);
CREATE INDEX friend_invites_with_user_id_idx ON friend_invites (with_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friend_invites;
-- +goose StatementEnd
