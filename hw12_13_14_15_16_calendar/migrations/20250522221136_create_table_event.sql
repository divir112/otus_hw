-- +goose Up
-- +goose StatementBegin
CREATE TABLE event(
    id BIGSERIAL PRIMARY KEY,
    title TEXT,
    start_time timestamp with time zone,
    end_time timestamp with time zone,
    description TEXT,
    user_id INTEGER,
    ping_before INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
