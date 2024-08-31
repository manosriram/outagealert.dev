-- +goose Up
-- +goose StatementBegin
ALTER TABLE ping ADD COLUMN status integer DEFAULT 200 NOT NULL, ADD COLUMN metadata jsonb;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ping DROP COLUMN status, DROP COLUMN metadata;
-- +goose StatementEnd
