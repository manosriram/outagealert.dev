-- +goose Up
-- +goose StatementBegin
ALTER TABLE project ADD COLUMN is_active boolean DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project DROP COLUMN is_active;
-- +goose StatementEnd
