-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ADD COLUMN period_text varchar(16) DEFAULT 'minutes';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor DROP COLUMN period_text;
-- +goose StatementEnd
