-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ADD COLUMN last_paused_at timestamp;
ALTER TABLE monitor ADD COLUMN last_resumed_at timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor DROP COLUMN last_paused_at;
ALTER TABLE monitor DROP COLUMN last_resumed_at;
-- +goose StatementEnd
