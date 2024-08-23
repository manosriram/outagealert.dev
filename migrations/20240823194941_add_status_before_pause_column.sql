-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitor ADD COLUMN status_before_pause varchar(64) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitor DROP COLUMN status_before_pause;
-- +goose StatementEnd
