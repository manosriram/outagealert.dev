-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_orders ADD COLUMN order_metadata json;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_orders DROP COLUMN order_metadata;
-- +goose StatementEnd
