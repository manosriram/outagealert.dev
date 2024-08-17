-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN project_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN project_id DROP NOT NULL;
-- +goose StatementEnd
