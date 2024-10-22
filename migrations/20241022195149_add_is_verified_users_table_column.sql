-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD is_verified boolean DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_verified;
-- +goose StatementEnd
