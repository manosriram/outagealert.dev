-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN created_at drop default;

ALTER TABLE monitor ALTER COLUMN created_at set default CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
