-- +goose Up
-- +goose StatementBegin
ALTER TABLE alert_integration ALTER COLUMN is_active SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE alert_integration ALTER COLUMN is_active DROP NOT NULL;
-- +goose StatementEnd
