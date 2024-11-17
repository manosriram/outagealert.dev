-- +goose Up
-- +goose StatementBegin
create unique index if not exists users_email_idx on users(email);
create index if not exists monitor_user_email_idx on monitor(user_email);
create index if not exists project_user_email_idx on project(user_email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index users_email_idx;
drop index monitor_user_email_idx;
drop index project_user_email_idx;
-- +goose StatementEnd
