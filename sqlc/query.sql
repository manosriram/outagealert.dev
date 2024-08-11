-- name: AllUsers :many
SELECT * FROM USERS;

-- name: CreateUser :one
INSERT INTO USERS(email, password) VALUES($1, $2) RETURNING *;

-- name: GetUserUsingEmail :one
SELECT * FROM USERS WHERE email = $1;

-- name: UpdateUserOtp :exec
UPDATE USERS SET otp = $1 WHERE email = $2;
