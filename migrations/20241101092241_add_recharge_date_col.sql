-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN recharge_date timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN recharge_date;
-- +goose StatementEnd
