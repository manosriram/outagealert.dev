-- +goose Up
-- +goose StatementBegin
CREATE TABLE plan (
		id serial,
		name varchar(64) UNIQUE,
		price int,
		validity int,
		PRIMARY KEY (id)
);
ALTER TABLE users ADD plan varchar(64);
ALTER TABLE users ADD CONSTRAINT userplan FOREIGN KEY (plan) REFERENCES plan(name);
ALTER TABLE users ALTER COLUMN plan SET DEFAULT 'free';
INSERT INTO plan(name, price, validity) VALUES('free', 0, 1000000000);
INSERT INTO plan(name, price, validity) VALUES('hobbyist', 3, 30);
INSERT INTO plan(name, price, validity) VALUES('pro', 10, 30);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE plan;
-- +goose StatementEnd
