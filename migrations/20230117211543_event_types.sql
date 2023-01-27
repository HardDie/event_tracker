-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_types (
    id         INTEGER   PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER   NOT NULL REFERENCES users(id),
    event_type TEXT      NOT NULL,
    is_visible BOOL      NOT NULL DEFAULT (true),
    created_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at TIMESTAMP
);
CREATE INDEX event_types_id_idx ON event_types (id);
CREATE INDEX event_types_user_id_idx ON event_types (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_types;
-- +goose StatementEnd
