-- +goose Up
-- +goose StatementBegin
alter table users add column otp integer;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users drop column otp;
-- +goose StatementEnd
