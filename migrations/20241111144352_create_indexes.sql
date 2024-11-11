-- +goose Up
-- +goose StatementBegin
create unique index users_email_idx on users(email);
create index monitor_user_email_idx on monitor(user_email);
create index project_user_email_idx on project(user_email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index users_email_idx;
drop index monitor_user_email_idx;
drop index project_user_email_idx;
-- +goose StatementEnd
