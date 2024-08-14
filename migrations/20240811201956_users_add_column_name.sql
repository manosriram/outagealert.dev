-- +goose Up
-- +goose StatementBegin
alter table users add column name varchar(32);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users drop column name;
-- +goose StatementEnd
