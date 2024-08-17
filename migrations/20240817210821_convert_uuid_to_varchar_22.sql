-- +goose Up
-- +goose StatementBegin
ALTER TABLE outagealertio.public.monitor DROP constraint monitor_pkey CASCADE;
alter table outagealertio.public.monitor alter column id drop default;
ALTER TABLE outagealertio.public.monitor ALTER COLUMN id TYPE varchar(22);
ALTER TABLE outagealertio.public.project DROP constraint project_pkey CASCADE;
alter table outagealertio.public.project alter column id drop default;
ALTER TABLE outagealertio.public.project ALTER COLUMN id TYPE varchar(22);
ALTER TABLE outagealertio.public.ping ALTER COLUMN id TYPE varchar(22);
alter table outagealertio.public.ping alter column id drop default;
ALTER TABLE outagealertio.public.project add PRIMARY KEY(id);
ALTER TABLE outagealertio.public.monitor add PRIMARY KEY(id);
alter table outagealertio.public.monitor alter column project_id type varchar(22);
alter table outagealertio.public.monitor alter column created_at type timestamptz;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
