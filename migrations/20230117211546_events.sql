-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id         SERIAL    PRIMARY KEY,
    user_id    INT       NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    type_id    INT       NOT NULL REFERENCES event_types(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    date       TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    deleted_at TIMESTAMP
);
CREATE INDEX events_user_id_date_idx ON events (user_id, date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
