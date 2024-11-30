-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS slack_users (
		user_email varchar(64) REFERENCES users(email) NOT NULL,
		channel_url varchar(512),
		channel_id varchar(64),
		channel_name varchar(256),
		configuration_url varchar(512),
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
CREATE TRIGGER updated_at BEFORE UPDATE ON slack_users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS slack_users;
-- +goose StatementEnd
