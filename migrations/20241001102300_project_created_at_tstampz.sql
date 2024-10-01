-- +goose Up
-- +goose StatementBegin
ALTER TABLE project ALTER COLUMN created_at TYPE timestamptz;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
