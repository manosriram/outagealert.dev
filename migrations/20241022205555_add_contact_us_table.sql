-- +goose Up
-- +goose StatementBegin
CREATE TABLE contact_us (
		name varchar(164),
		email varchar(64),
		message text,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp default CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table contact_us;
-- +goose StatementEnd
