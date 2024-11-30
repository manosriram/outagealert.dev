-- +goose Up
-- +goose StatementBegin
ALTER TABLE slack_users DROP COLUMN user_email;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
