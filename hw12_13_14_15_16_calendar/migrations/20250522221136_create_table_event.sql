-- +goose Up
-- +goose StatementBegin
CREATE TABLE event(
    ID BIGSERIAL PRIMARY KEY,
    Header TEXT,
    Date timestamp with time zone,
    DateEnd timestamp with time zone,
    Description TEXT,
    Owner TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
