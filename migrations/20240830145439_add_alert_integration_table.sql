-- +goose Up
-- +goose StatementBegin
CREATE TYPE ALERT_TYPE as ENUM ('email', 'slack', 'webhook');
CREATE TABLE IF NOT EXISTS alert_integration (
		id varchar(22),
		monitor_id varchar(22) REFERENCES monitor(id) NOT NULL,
		is_active boolean DEFAULT true,
		alert_type ALERT_TYPE NOT NULL,
		alert_target varchar(512),
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
CREATE TRIGGER updated_at BEFORE UPDATE ON alert_integration FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE ALERT_TYPE CASCADE;
DROP TABLE IF EXISTS alert_integration;
-- +goose StatementEnd
