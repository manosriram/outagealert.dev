-- +goose Up
-- +goose StatementBegin
create index ping_monitor_idx on ping(monitor_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index ping_monitor_idx;
-- +goose StatementEnd
