-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id         INTEGER   PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER   NOT NULL REFERENCES users(id),
    type_id    INTEGER   NOT NULL REFERENCES event_types(id),
    date       TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at TIMESTAMP
);
CREATE INDEX events_id_idx ON events (id);
CREATE INDEX events_user_id_date_idx ON events (user_id, date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
