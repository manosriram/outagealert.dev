-- +goose Up
-- +goose StatementBegin
ALTER TABLE alert_integration ADD COLUMN email_alert_sent boolean DEFAULT false;
ALTER TABLE alert_integration ADD COLUMN slack_alert_sent boolean DEFAULT false;
ALTER TABLE alert_integration ADD COLUMN webhook_alert_sent boolean DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE alert_integration DROP COLUMN email_alert_sent;
ALTER TABLE alert_integration DROP COLUMN slack_alert_sent;
ALTER TABLE alert_integration DROP COLUMN webhook_alert_sent;
-- +goose StatementEnd
