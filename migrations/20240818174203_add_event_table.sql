-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event (
		id varchar(22) primary key,
		monitor_id varchar(22) references monitor(id) NOT NULL,
		from_status varchar(64),
		to_status varchar(64),
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
