-- +goose Up
-- +goose StatementBegin
ALTER TABLE USERS ALTER COLUMN otp TYPE varchar(32);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE USERS ALTER COLUMN otp int;
-- +goose StatementEnd
