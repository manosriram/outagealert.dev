-- +goose Up
-- +goose StatementBegin
SELECT 'CREATE DATABASE outagealert' 
WHERE NOT EXISTS (
    SELECT FROM pg_database WHERE datname = 'outagealert'
)\gexec
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE DATABASE outagealert;
-- +goose StatementEnd
