-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
CREATE DATABASE calendar WITH Owner = postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE IF EXISTS calendar;
-- +goose StatementEnd
