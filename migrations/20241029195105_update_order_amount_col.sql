-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_orders DROP COLUMN order_amount; 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_orders ADD COLUMN order_amount decimal; 
-- +goose StatementEnd
