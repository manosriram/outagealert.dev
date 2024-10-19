-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.monitor DROP constraint monitor_pkey CASCADE;
alter table public.monitor alter column id drop default;
ALTER TABLE public.monitor ALTER COLUMN id TYPE varchar(22);
ALTER TABLE public.project DROP constraint project_pkey CASCADE;
alter table public.project alter column id drop default;
ALTER TABLE public.project ALTER COLUMN id TYPE varchar(22);
ALTER TABLE public.ping ALTER COLUMN id TYPE varchar(22);
alter table public.ping alter column id drop default;
ALTER TABLE public.project add PRIMARY KEY(id);
ALTER TABLE public.monitor add PRIMARY KEY(id);
alter table public.monitor alter column project_id type varchar(22);
alter table public.monitor alter column created_at type timestamptz;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
