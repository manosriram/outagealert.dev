-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN last_ping SET DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN last_ping DROP NULL;
-- +goose StatementEnd
