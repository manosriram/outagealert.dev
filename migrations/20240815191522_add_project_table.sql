-- +goose Up
-- +goose StatementBegin
--  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--  CREATE EXTENSION pgcrypto;
CREATE TABLE IF NOT EXISTS project (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name varchar(64) NOT NULL,
		user_email varchar(64) REFERENCES users(email) NOT NULL,
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
CREATE TRIGGER updated_at BEFORE UPDATE ON project FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project;
-- +goose StatementEnd
