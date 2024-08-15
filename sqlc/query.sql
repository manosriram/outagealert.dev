-- name: AllUsers :many
SELECT * FROM USERS;

-- name: Create :one
INSERT INTO USERS(name, email, password) VALUES($1, $2, $3) RETURNING *;

-- name: GetUserUsingEmail :one
SELECT * FROM USERS WHERE email = $1;

-- name: GetUserUsingOtp :one
SELECT * FROM USERS WHERE otp = $1;

-- name: UpdateUserOtp :exec
UPDATE USERS SET otp = $1 WHERE email = $2;

-- name: ResetUserPassword :exec
UPDATE USERS SET password = $1, otp = '' WHERE email = $2;

-- name: CreateMonitor :exec
INSERT INTO monitor(name, period, grace_period, user_id, project_id, ping_url, type) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetMonitorById :one
SELECT * FROM monitor WHERE id = $1;

-- name: UpdateUserMonitorName :exec
UPDATE monitor SET name = $1 WHERE user_id = $2;

-- name: UpdateUserMonitorSchedule :exec
UPDATE monitor SET period = $1, grace_period = $2 WHERE user_id = $3;

-- name: GetUserMonitors :many
SELECT * FROM monitor WHERE user_id = $1;

-- name: CreateProject :exec
INSERT INTO project(name, user_id) VALUES($1, $2) RETURNING *;

-- name: UpdateUserProjectName :exec
UPDATE project SET name = $1 WHERE user_id = $2;

-- name: GetProjectById :one
SELECT * FROM project WHERE id = $1;

-- name: GetUserProjects :many
SELECT * FROM project WHERE user_id = $1;
