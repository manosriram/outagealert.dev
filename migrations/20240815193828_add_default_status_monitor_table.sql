-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN status SET DEFAULT 'up';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN status DROP DEFAULT;
-- +goose StatementEnd
