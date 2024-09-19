-- +goose Up
-- +goose StatementBegin
--  ALTER TABLE alert_integration DROP CONSTRAINT alert_integration_pkey;
ALTER TABLE alert_integration ADD PRIMARY KEY (monitor_id, alert_type);
-- +goose StatementEnd
