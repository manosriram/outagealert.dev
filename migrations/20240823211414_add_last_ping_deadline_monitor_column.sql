-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ADD COLUMN total_pause_time integer DEFAULT 0; -- in seconds
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor DROP COLUMN total_pause_time; -- in seconds
-- +goose StatementEnd
