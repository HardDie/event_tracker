-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS friend_invites (
    id           SERIAL    PRIMARY KEY,
    user_id      INT       NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    with_user_id INT       NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at   TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at   TIMESTAMP
);
CREATE INDEX friend_invites_user_id_idx ON friend_invites (user_id);
CREATE INDEX friend_invites_with_user_id_idx ON friend_invites (with_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friend_invites;
-- +goose StatementEnd
