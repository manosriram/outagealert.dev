-- +goose Up
-- +goose StatementBegin
ALTER TABLE project ADD COLUMN visibility varchar(32) default 'public';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project DROP COLUMN visibility;
-- +goose StatementEnd
