-- +goose Up
-- +goose StatementBegin
--  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS monitor (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name varchar(64) NOT NULL,
		period integer NOT NULL DEFAULT 600,
		grace_period integer NOT NULL DEFAULT 300,
		user_email varchar(64) REFERENCES users(email) NOT NULL,
		project_id UUID NOT NULL,
		ping_url varchar(512) NOT NULL,
		status varchar(64) NULL,
		type varchar(64) NULL,
		last_ping timestamp DEFAULT CURRENT_TIMESTAMP NULL,
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

CREATE TRIGGER updated_at BEFORE UPDATE ON monitor FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS monitor;
-- +goose StatementEnd
