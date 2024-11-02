-- name: AllUsers :many
SELECT * FROM USERS;

-- name: Create :one
INSERT INTO USERS(name, email, password, magic_token, uuid) VALUES($1, $2, $3, $4, $5) RETURNING *;

-- name: MarkUserVerified :exec
UPDATE USERS SET is_verified = true WHERE email = $1;

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

-- name: UpdateMonitor :exec
UPDATE monitor SET name = $1, period = $2, grace_period = $3 WHERE id = $4 AND user_email = $5;

-- name: DeleteMonitor :exec
UPDATE monitor SET is_active=false WHERE id = $1 AND user_email = $2;

-- name: GetMonitorById :one
SELECT * FROM monitor where id = $1 AND is_active=true;

-- name: GetMonitorWithEventsById :many
SELECT * FROM monitor m JOIN event e ON m.id = e.monitor_id AND m.id = $1 AND m.is_active = true;

-- name: GetMonitorByPingUrl :one
SELECT * FROM monitor m WHERE ping_url = $1 AND m.is_active = true;

-- name: UpdateUserMonitorName :exec
UPDATE monitor SET name = $1 WHERE user_email = $2;

-- name: UpdateUserMonitorSchedule :exec
UPDATE monitor SET period = $1, grace_period = $2 WHERE user_email = $3;

-- name: UpdateMonitorLastPing :exec
UPDATE monitor SET last_ping = $1, status='up' WHERE id = $2;

-- name: UpdateMonitorStatus :exec
UPDATE monitor SET status = $1 WHERE id = $2;

-- name: PauseMonitor :one
UPDATE monitor SET status = $1, status_before_pause = $2, last_paused_at = $3 WHERE id = $4 RETURNING *;

-- name: ResumeMonitor :one
UPDATE monitor m SET status = m.status_before_pause, status_before_pause = '', last_resumed_at = $1, total_pause_time = $2 WHERE id = $3 RETURNING *;

-- name: GetAllMonitorIDs :many
SELECT id, period, grace_period from monitor where is_active = true;

-- name: GetUserMonitors :many
SELECT * FROM monitor WHERE user_email = $1;

-- name: CreateProject :one
INSERT INTO project(id, name, user_email, visibility) VALUES($1, $2, $3, $4) RETURNING *;

-- name: UpdateUserProjectName :exec
UPDATE project SET name = $1 WHERE user_email = $2;

-- name: GetProjectById :one
SELECT * FROM project WHERE id = $1;

-- name: GetUserProjects :many
SELECT p.*, COUNT(m.id) AS monitor_count FROM project p LEFT JOIN monitor m ON p.id = m.project_id AND m.is_active = true WHERE p.user_email = $1 GROUP BY p.id ORDER BY p.created_at DESC;

-- name: GetProjectMonitors :many
SELECT * FROM monitor WHERE project_id = $1 AND is_active=true ORDER BY created_at DESC;

-- name: CreatePing :exec
INSERT INTO ping(id, monitor_id, status, metadata) VALUES($1, $2, $3, $4) RETURNING *;

-- name: GetMonitorPings :many
SELECT * FROM ping where monitor_id = $1;

-- name: CreateEvent :exec
INSERT INTO event(id, monitor_id, from_status, to_status, created_at) VALUES($1, $2, $3, $4, $5) RETURNING *;

-- name: GetEventById :many
SELECT * FROM event WHERE id = $1;

-- name: GetEventsByMonitorId :many
SELECT * FROM event where monitor_id = $1;

-- name: GetEventsByMonitorIdPaginated :many
SELECT * FROM event where monitor_id = $1 ORDER BY created_at DESC LIMIT 25 OFFSET $2;

-- name: GetPingsByMonitorIdPaginated :many
SELECT * FROM ping where monitor_id = $1 ORDER BY created_at DESC LIMIT 25 OFFSET $2;

-- name: GetLastToStatusUpMonitorEvent :one
SELECT * FROM event where monitor_id = $1 AND to_status='up' AND from_status != 'up' order by created_at desc;

-- name: GetLatestNonPausedMonitorEvent :one
SELECT * FROM event where monitor_id = $1 AND to_status != 'paused' order by created_at desc;

-- name: GetLastToPausedMonitorEvent :one
SELECT * FROM event where monitor_id = $1 AND to_status='paused' order by created_at desc;

-- name: GetLatestMonitorEventByToStatus :one
SELECT * FROM event where monitor_id = $1 AND to_status=$2 order by created_at desc;

-- name: GetNumberOfMonitorIncidents :one
SELECT count(*) FROM event where monitor_id = $1 AND (from_status='grace_period' or from_status='up') AND to_status='down';

-- name: UpdateMonitorTotalPauseTime :exec
UPDATE monitor set total_pause_time = $1 where id = $2;

-- name: TotalMonitorPings :one
SELECT COUNT(*) as ping_count FROM ping where monitor_id = $1;

-- name: TotalMonitorEvents :one
SELECT COUNT(*) as event_count FROM event WHERE monitor_id = $1;

-- name: GetMonitorActivityPaginated :many
SELECT id, from_status, to_status, created_at, updated_at, source, status, metadata
FROM (
    SELECT id, from_status, to_status, 
           created_at AT TIME ZONE 'UTC' AS created_at, 
           updated_at AT TIME ZONE 'UTC' AS updated_at, 
           'event' AS source,
           200 AS status, 
           NULL::jsonb AS metadata
    FROM event e
    WHERE e.monitor_id = $1
    UNION ALL
    SELECT id, 'active' AS from_status, 'active' AS to_status, 
           created_at, updated_at, 
           'ping' AS source,
           status, 
           metadata::jsonb
    FROM ping p
    WHERE p.monitor_id = $1
) AS combined
ORDER BY created_at DESC
LIMIT 25 OFFSET $2;

-- name: InitMonitorIntegrations :exec
INSERT INTO alert_integration(id, monitor_id, alert_type, is_active) VALUES($1, $2, $3, $4);

-- name: GetMonitorIntegrations :many
SELECT * FROM alert_integration WHERE monitor_id = $1 ORDER BY created_at;

-- name: GetMonitorIntegration :one
SELECT * FROM alert_integration WHERE monitor_id = $1 AND alert_type = $2 ORDER BY created_at;

-- name: UpdateEmailAlertIntegration :exec
UPDATE alert_integration set is_active = $1 WHERE monitor_id = $2 AND alert_type = 'email';

-- name: UpdateSlackAlertIntegration :exec
UPDATE alert_integration set is_active = $1 WHERE monitor_id = $2 AND alert_type = 'slack';

-- name: UpdateWebhookAlertIntegration :exec
UPDATE alert_integration set alert_target = $1, is_active = $2 WHERE monitor_id = $3 AND alert_type = 'webhook';

-- name: GetMonitorAlertIntegration :one
SELECT * FROM alert_integration WHERE monitor_id = $1 AND alert_type = $2;

-- name: UpdateEmailAlertSentFlag :exec
UPDATE alert_integration set email_alert_sent = $1 WHERE monitor_id = $2;

-- name: UpdateSlackAlertSentFlag :exec
UPDATE alert_integration set slack_alert_sent = $1 WHERE monitor_id = $2;

-- name: UpdateWebhookAlertSentFlag :exec
UPDATE alert_integration set webhook_alert_sent = $1 WHERE monitor_id = $2;

-- name: UpdateAlertSentFlag :exec
UPDATE alert_integration set email_alert_sent = $1, slack_alert_sent = $2, webhook_alert_sent = $3 WHERE monitor_id = $4;

-- name: UserMonitorCount :one
select count(m.id) from monitor m where user_email = $1 and m.is_active = true;

-- name: UpdateUserMagicToken :exec
UPDATE users set magic_token = $1 WHERE id = $2;

-- name: AddContactFormEntry :exec
INSERT INTO contact_us(name, email, message) VALUES($1, $2, $3);

-- name: CreateOrder :exec
INSERT INTO user_orders(
		order_id,
		user_email,
		order_status,
		order_payment_session_id,
		plan,
		order_expiry_time,
		order_currency
) VALUES($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateUserPlan :exec
UPDATE users SET plan = $1, recharge_date = NOW() WHERE email = $2;

-- name: UpdateOrderStatusAndMetadata :exec
UPDATE user_orders SET order_status = $1, order_metadata = $2 WHERE order_id = $3;

-- name: GetOrderByOrderId :one
SELECT order_id, user_email, order_status, order_payment_session_id, plan, created_at FROM user_orders WHERE order_id = $1;
