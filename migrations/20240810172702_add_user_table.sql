-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS USERS (
		id serial,
		email varchar(64),
		password varchar(256),
		is_active boolean default true,
		last_login timestamp,
		created_at timestamp default current_timestamp,
		updated_at timestamp default current_timestamp,
		primary key (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS USERS;
-- +goose StatementEnd
