-- +goose Up
-- +goose StatementBegin
create index ping_monitor_created_at_idx on ping(monitor_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index ping_monitor_created_at_idx;
-- +goose StatementEnd
