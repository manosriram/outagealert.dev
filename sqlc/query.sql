-- name: AllUsers :many
SELECT * FROM USERS;

-- name: CreateUser :one
INSERT INTO USERS(email, password) VALUES($1, $2) RETURNING *;
