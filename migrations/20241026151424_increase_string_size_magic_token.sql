-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN magic_token TYPE varchar(512);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN magic_token TYPE varchar(22);
-- +goose StatementEnd
