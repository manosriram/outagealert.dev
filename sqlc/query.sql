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

-- name: CreateMonitor :one
INSERT INTO monitor(id, name, period, grace_period, user_email, project_id, ping_url, type) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetMonitorById :one
SELECT * FROM monitor where id = $1;

-- name: GetMonitorWithEventsById :many
SELECT * FROM monitor m JOIN event e ON m.id = e.monitor_id AND m.id = $1;

-- name: GetMonitorByPingUrl :one
SELECT * FROM monitor WHERE ping_url = $1;

-- name: UpdateUserMonitorName :exec
UPDATE monitor SET name = $1 WHERE user_email = $2;

-- name: UpdateUserMonitorSchedule :exec
UPDATE monitor SET period = $1, grace_period = $2 WHERE user_email = $3;

-- name: UpdateMonitorLastPing :exec
UPDATE monitor SET last_ping = $1, status='up' WHERE id = $2;

-- name: UpdateMonitorStatus :exec
UPDATE monitor SET status = $1 WHERE id = $2;

-- name: GetAllMonitorIDs :many
SELECT id, period, grace_period from monitor;

-- name: GetUserMonitors :many
SELECT * FROM monitor WHERE user_email = $1;

-- name: CreateProject :one
INSERT INTO project(id, name, user_email, visibility) VALUES($1, $2, $3, $4) RETURNING *;

-- name: UpdateUserProjectName :exec
UPDATE project SET name = $1 WHERE user_email = $2;

-- name: GetProjectById :one
SELECT * FROM project WHERE id = $1;

-- name: GetUserProjects :many
SELECT * FROM project WHERE user_email = $1;

-- name: GetProjectMonitors :many
SELECT * FROM monitor WHERE project_id = $1;

-- name: CreatePing :exec
INSERT INTO ping(id, monitor_id) VALUES($1, $2) RETURNING *;

-- name: GetMonitorPings :many
SELECT * FROM ping where monitor_id = $1;

-- name: CreateEvent :exec
INSERT INTO event(id, monitor_id, from_status, to_status) VALUES($1, $2, $3, $4) RETURNING *;

-- name: GetEventById :many
SELECT * FROM event WHERE id = $1;

-- name: GetEventsByMonitorId :many
SELECT * FROM event where monitor_id = $1;

-- name: GetEventsByMonitorIdPaginated :many
SELECT * FROM event where monitor_id = $1 ORDER BY created_at DESC LIMIT 25 OFFSET $2;
