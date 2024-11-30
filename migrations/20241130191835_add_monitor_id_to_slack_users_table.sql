-- +goose Up
-- +goose StatementBegin
ALTER TABLE slack_users ADD COLUMN monitor_id varchar(22);
ALTER TABLE slack_users DROP CONSTRAINT pk_slack_users_email;
ALTER TABLE slack_users ADD PRIMARY KEY(monitor_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE slack_users DROP COLUMN monitor_id; -- not possible since it is a primary key
-- +goose StatementEnd
