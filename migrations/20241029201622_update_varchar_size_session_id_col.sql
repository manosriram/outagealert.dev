-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_orders ALTER COLUMN order_payment_session_id TYPE varchar(256);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_orders ALTER COLUMN order_payment_session_id TYPE varchar(64);
-- +goose StatementEnd
