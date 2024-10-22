-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD magic_token varchar(22);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP magic_token;
-- +goose StatementEnd
