-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN uuid varchar(22);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN uuid;
-- +goose StatementEnd
