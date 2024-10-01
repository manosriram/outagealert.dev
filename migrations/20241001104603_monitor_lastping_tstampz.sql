-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN last_ping TYPE timestamptz;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
