-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS friends (
    id           SERIAL    PRIMARY KEY,
    user_id      INT       NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    with_user_id INT       NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at   TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at   TIMESTAMP
);
CREATE INDEX friends_user_id_idx ON friends (user_id);
CREATE INDEX friends_with_user_id_idx ON friends (with_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friends;
-- +goose StatementEnd
