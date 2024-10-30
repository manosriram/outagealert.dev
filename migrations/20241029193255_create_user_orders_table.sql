-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_orders (
		order_id varchar(256) PRIMARY KEY,
		user_email varchar(64) REFERENCES users(email) NOT NULL,
		order_status varchar(16) NOT NULL,
		order_payment_session_id varchar(64),
		plan varchar(22),
		order_expiry_time timestamp,
		order_currency varchar(16),
		order_amount decimal,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL
);
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
		NEW.updated_at = now();
		RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER updated_at BEFORE UPDATE ON user_orders FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_orders;
-- +goose StatementEnd
