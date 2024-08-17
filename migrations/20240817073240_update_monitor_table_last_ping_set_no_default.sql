-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN last_ping DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN last_ping SET DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd
