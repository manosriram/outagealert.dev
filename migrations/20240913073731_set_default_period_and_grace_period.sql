-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN period SET DEFAULT 25;
ALTER TABLE monitor ALTER COLUMN grace_period SET DEFAULT 25;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor ALTER COLUMN period DROP DEFAULT;
ALTER TABLE monitor ALTER COLUMN grace_period DROP DEFAULT;
-- +goose StatementEnd
