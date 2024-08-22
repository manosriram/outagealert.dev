-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ADD COLUMN is_active boolean DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor DROP COLUMN is_active;
-- +goose StatementEnd
