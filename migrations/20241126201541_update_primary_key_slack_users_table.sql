-- +goose Up
-- +goose StatementBegin
ALTER TABLE slack_users
ADD CONSTRAINT PK_slack_users_email PRIMARY KEY (user_email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE slack_users
DROP PRIMARY KEY;
-- +goose StatementEnd
