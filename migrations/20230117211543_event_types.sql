-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_types (
    id         SERIAL    PRIMARY KEY,
    user_id    INT       NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    event_type TEXT      NOT NULL,
    is_visible BOOL      NOT NULL DEFAULT (true),
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at TIMESTAMP
);
CREATE INDEX event_types_user_id_idx ON event_types (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_types;
-- +goose StatementEnd
