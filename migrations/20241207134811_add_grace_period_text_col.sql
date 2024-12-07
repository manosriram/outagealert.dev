-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ADD COLUMN grace_period_text varchar(16) DEFAULT 'minutes';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor DROP COLUMN grace_period_text;
-- +goose StatementEnd
