-- +goose Up
-- +goose StatementBegin
ALTER TABLE ping ALTER COLUMN monitor_id TYPE varchar(22);
ALTER TABLE ping ADD CONSTRAINT fk_monitor_id FOREIGN KEY (monitor_id) REFERENCES monitor(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
