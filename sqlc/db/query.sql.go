// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const allUsers = `-- name: AllUsers :many
SELECT id, name, email, password, is_active, otp, last_login, created_at, updated_at FROM USERS
`

func (q *Queries) AllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, allUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.IsActive,
			&i.Otp,
			&i.LastLogin,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const create = `-- name: Create :one
INSERT INTO USERS(name, email, password) VALUES($1, $2, $3) RETURNING id, name, email, password, is_active, otp, last_login, created_at, updated_at
`

type CreateParams struct {
	Name     *string
	Email    string
	Password string
}

func (q *Queries) Create(ctx context.Context, arg CreateParams) (User, error) {
	row := q.db.QueryRow(ctx, create, arg.Name, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.Otp,
		&i.LastLogin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createEvent = `-- name: CreateEvent :exec
INSERT INTO event(id, monitor_id, from_status, to_status, created_at) VALUES($1, $2, $3, $4, $5) RETURNING id, monitor_id, from_status, to_status, created_at, updated_at
`

type CreateEventParams struct {
	ID         string
	MonitorID  string
	FromStatus string
	ToStatus   string
	CreatedAt  pgtype.Timestamp
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) error {
	_, err := q.db.Exec(ctx, createEvent,
		arg.ID,
		arg.MonitorID,
		arg.FromStatus,
		arg.ToStatus,
		arg.CreatedAt,
	)
	return err
}

const createMonitor = `-- name: CreateMonitor :one
INSERT INTO monitor(id, name, period, grace_period, user_email, project_id, ping_url, type) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, name, period, grace_period, user_email, project_id, ping_url, status, status_before_pause, is_active, type, total_pause_time, last_ping, last_paused_at, last_resumed_at, created_at, updated_at
`

type CreateMonitorParams struct {
	ID          string
	Name        string
	Period      int32
	GracePeriod int32
	UserEmail   string
	ProjectID   string
	PingUrl     string
	Type        string
}

func (q *Queries) CreateMonitor(ctx context.Context, arg CreateMonitorParams) (Monitor, error) {
	row := q.db.QueryRow(ctx, createMonitor,
		arg.ID,
		arg.Name,
		arg.Period,
		arg.GracePeriod,
		arg.UserEmail,
		arg.ProjectID,
		arg.PingUrl,
		arg.Type,
	)
	var i Monitor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Period,
		&i.GracePeriod,
		&i.UserEmail,
		&i.ProjectID,
		&i.PingUrl,
		&i.Status,
		&i.StatusBeforePause,
		&i.IsActive,
		&i.Type,
		&i.TotalPauseTime,
		&i.LastPing,
		&i.LastPausedAt,
		&i.LastResumedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createPing = `-- name: CreatePing :exec
INSERT INTO ping(id, monitor_id) VALUES($1, $2) RETURNING id, monitor_id, created_at, updated_at
`

type CreatePingParams struct {
	ID        string
	MonitorID string
}

func (q *Queries) CreatePing(ctx context.Context, arg CreatePingParams) error {
	_, err := q.db.Exec(ctx, createPing, arg.ID, arg.MonitorID)
	return err
}

const createProject = `-- name: CreateProject :one
INSERT INTO project(id, name, user_email, visibility) VALUES($1, $2, $3, $4) RETURNING id, name, visibility, user_email, created_at, updated_at
`

type CreateProjectParams struct {
	ID         string
	Name       string
	UserEmail  string
	Visibility string
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRow(ctx, createProject,
		arg.ID,
		arg.Name,
		arg.UserEmail,
		arg.Visibility,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Visibility,
		&i.UserEmail,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteMonitor = `-- name: DeleteMonitor :exec
UPDATE monitor SET is_active=false WHERE id = $1 AND user_email = $2
`

type DeleteMonitorParams struct {
	ID        string
	UserEmail string
}

func (q *Queries) DeleteMonitor(ctx context.Context, arg DeleteMonitorParams) error {
	_, err := q.db.Exec(ctx, deleteMonitor, arg.ID, arg.UserEmail)
	return err
}

const getAllMonitorIDs = `-- name: GetAllMonitorIDs :many
SELECT id, period, grace_period from monitor
`

type GetAllMonitorIDsRow struct {
	ID          string
	Period      int32
	GracePeriod int32
}

func (q *Queries) GetAllMonitorIDs(ctx context.Context) ([]GetAllMonitorIDsRow, error) {
	rows, err := q.db.Query(ctx, getAllMonitorIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllMonitorIDsRow
	for rows.Next() {
		var i GetAllMonitorIDsRow
		if err := rows.Scan(&i.ID, &i.Period, &i.GracePeriod); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventById = `-- name: GetEventById :many
SELECT id, monitor_id, from_status, to_status, created_at, updated_at FROM event WHERE id = $1
`

func (q *Queries) GetEventById(ctx context.Context, id string) ([]Event, error) {
	rows, err := q.db.Query(ctx, getEventById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.MonitorID,
			&i.FromStatus,
			&i.ToStatus,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventsByMonitorId = `-- name: GetEventsByMonitorId :many
SELECT id, monitor_id, from_status, to_status, created_at, updated_at FROM event where monitor_id = $1
`

func (q *Queries) GetEventsByMonitorId(ctx context.Context, monitorID string) ([]Event, error) {
	rows, err := q.db.Query(ctx, getEventsByMonitorId, monitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.MonitorID,
			&i.FromStatus,
			&i.ToStatus,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventsByMonitorIdPaginated = `-- name: GetEventsByMonitorIdPaginated :many
SELECT id, monitor_id, from_status, to_status, created_at, updated_at FROM event where monitor_id = $1 ORDER BY created_at DESC LIMIT 25 OFFSET $2
`

type GetEventsByMonitorIdPaginatedParams struct {
	MonitorID string
	Offset    int32
}

func (q *Queries) GetEventsByMonitorIdPaginated(ctx context.Context, arg GetEventsByMonitorIdPaginatedParams) ([]Event, error) {
	rows, err := q.db.Query(ctx, getEventsByMonitorIdPaginated, arg.MonitorID, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.MonitorID,
			&i.FromStatus,
			&i.ToStatus,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastToPausedMonitorEvent = `-- name: GetLastToPausedMonitorEvent :one
SELECT id, monitor_id, from_status, to_status, created_at, updated_at FROM event where monitor_id = $1 AND to_status='paused' order by created_at desc
`

func (q *Queries) GetLastToPausedMonitorEvent(ctx context.Context, monitorID string) (Event, error) {
	row := q.db.QueryRow(ctx, getLastToPausedMonitorEvent, monitorID)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.MonitorID,
		&i.FromStatus,
		&i.ToStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLastToStatusUpMonitorEvent = `-- name: GetLastToStatusUpMonitorEvent :one
SELECT id, monitor_id, from_status, to_status, created_at, updated_at FROM event where monitor_id = $1 AND to_status='up' AND from_status != 'up' order by created_at desc
`

func (q *Queries) GetLastToStatusUpMonitorEvent(ctx context.Context, monitorID string) (Event, error) {
	row := q.db.QueryRow(ctx, getLastToStatusUpMonitorEvent, monitorID)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.MonitorID,
		&i.FromStatus,
		&i.ToStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLatestMonitorEventByToStatus = `-- name: GetLatestMonitorEventByToStatus :one
SELECT id, monitor_id, from_status, to_status, created_at, updated_at FROM event where monitor_id = $1 AND to_status=$2 order by created_at desc
`

type GetLatestMonitorEventByToStatusParams struct {
	MonitorID string
	ToStatus  string
}

func (q *Queries) GetLatestMonitorEventByToStatus(ctx context.Context, arg GetLatestMonitorEventByToStatusParams) (Event, error) {
	row := q.db.QueryRow(ctx, getLatestMonitorEventByToStatus, arg.MonitorID, arg.ToStatus)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.MonitorID,
		&i.FromStatus,
		&i.ToStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLatestNonPausedMonitorEvent = `-- name: GetLatestNonPausedMonitorEvent :one
SELECT id, monitor_id, from_status, to_status, created_at, updated_at FROM event where monitor_id = $1 AND to_status != 'paused' order by created_at desc
`

func (q *Queries) GetLatestNonPausedMonitorEvent(ctx context.Context, monitorID string) (Event, error) {
	row := q.db.QueryRow(ctx, getLatestNonPausedMonitorEvent, monitorID)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.MonitorID,
		&i.FromStatus,
		&i.ToStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMonitorById = `-- name: GetMonitorById :one
SELECT id, name, period, grace_period, user_email, project_id, ping_url, status, status_before_pause, is_active, type, total_pause_time, last_ping, last_paused_at, last_resumed_at, created_at, updated_at FROM monitor where id = $1 AND is_active=true
`

func (q *Queries) GetMonitorById(ctx context.Context, id string) (Monitor, error) {
	row := q.db.QueryRow(ctx, getMonitorById, id)
	var i Monitor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Period,
		&i.GracePeriod,
		&i.UserEmail,
		&i.ProjectID,
		&i.PingUrl,
		&i.Status,
		&i.StatusBeforePause,
		&i.IsActive,
		&i.Type,
		&i.TotalPauseTime,
		&i.LastPing,
		&i.LastPausedAt,
		&i.LastResumedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMonitorByPingUrl = `-- name: GetMonitorByPingUrl :one
SELECT id, name, period, grace_period, user_email, project_id, ping_url, status, status_before_pause, is_active, type, total_pause_time, last_ping, last_paused_at, last_resumed_at, created_at, updated_at FROM monitor m WHERE ping_url = $1 AND m.is_active = true
`

func (q *Queries) GetMonitorByPingUrl(ctx context.Context, pingUrl string) (Monitor, error) {
	row := q.db.QueryRow(ctx, getMonitorByPingUrl, pingUrl)
	var i Monitor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Period,
		&i.GracePeriod,
		&i.UserEmail,
		&i.ProjectID,
		&i.PingUrl,
		&i.Status,
		&i.StatusBeforePause,
		&i.IsActive,
		&i.Type,
		&i.TotalPauseTime,
		&i.LastPing,
		&i.LastPausedAt,
		&i.LastResumedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMonitorPings = `-- name: GetMonitorPings :many
SELECT id, monitor_id, created_at, updated_at FROM ping where monitor_id = $1
`

func (q *Queries) GetMonitorPings(ctx context.Context, monitorID string) ([]Ping, error) {
	rows, err := q.db.Query(ctx, getMonitorPings, monitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Ping
	for rows.Next() {
		var i Ping
		if err := rows.Scan(
			&i.ID,
			&i.MonitorID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMonitorWithEventsById = `-- name: GetMonitorWithEventsById :many
SELECT m.id, name, period, grace_period, user_email, project_id, ping_url, status, status_before_pause, is_active, type, total_pause_time, last_ping, last_paused_at, last_resumed_at, m.created_at, m.updated_at, e.id, monitor_id, from_status, to_status, e.created_at, e.updated_at FROM monitor m JOIN event e ON m.id = e.monitor_id AND m.id = $1 AND m.is_active = true
`

type GetMonitorWithEventsByIdRow struct {
	ID                string
	Name              string
	Period            int32
	GracePeriod       int32
	UserEmail         string
	ProjectID         string
	PingUrl           string
	Status            string
	StatusBeforePause *string
	IsActive          *bool
	Type              string
	TotalPauseTime    *int32
	LastPing          pgtype.Timestamp
	LastPausedAt      pgtype.Timestamp
	LastResumedAt     pgtype.Timestamp
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
	ID_2              string
	MonitorID         string
	FromStatus        string
	ToStatus          string
	CreatedAt_2       pgtype.Timestamp
	UpdatedAt_2       pgtype.Timestamp
}

func (q *Queries) GetMonitorWithEventsById(ctx context.Context, id string) ([]GetMonitorWithEventsByIdRow, error) {
	rows, err := q.db.Query(ctx, getMonitorWithEventsById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMonitorWithEventsByIdRow
	for rows.Next() {
		var i GetMonitorWithEventsByIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Period,
			&i.GracePeriod,
			&i.UserEmail,
			&i.ProjectID,
			&i.PingUrl,
			&i.Status,
			&i.StatusBeforePause,
			&i.IsActive,
			&i.Type,
			&i.TotalPauseTime,
			&i.LastPing,
			&i.LastPausedAt,
			&i.LastResumedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.MonitorID,
			&i.FromStatus,
			&i.ToStatus,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNumberOfMonitorIncidents = `-- name: GetNumberOfMonitorIncidents :one
SELECT count(*) FROM event where monitor_id = $1 AND (from_status='grace_period' or from_status='up') AND to_status='down'
`

func (q *Queries) GetNumberOfMonitorIncidents(ctx context.Context, monitorID string) (int64, error) {
	row := q.db.QueryRow(ctx, getNumberOfMonitorIncidents, monitorID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getProjectById = `-- name: GetProjectById :one
SELECT id, name, visibility, user_email, created_at, updated_at FROM project WHERE id = $1
`

func (q *Queries) GetProjectById(ctx context.Context, id string) (Project, error) {
	row := q.db.QueryRow(ctx, getProjectById, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Visibility,
		&i.UserEmail,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProjectMonitors = `-- name: GetProjectMonitors :many
SELECT id, name, period, grace_period, user_email, project_id, ping_url, status, status_before_pause, is_active, type, total_pause_time, last_ping, last_paused_at, last_resumed_at, created_at, updated_at FROM monitor WHERE project_id = $1 AND is_active=true
`

func (q *Queries) GetProjectMonitors(ctx context.Context, projectID string) ([]Monitor, error) {
	rows, err := q.db.Query(ctx, getProjectMonitors, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Monitor
	for rows.Next() {
		var i Monitor
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Period,
			&i.GracePeriod,
			&i.UserEmail,
			&i.ProjectID,
			&i.PingUrl,
			&i.Status,
			&i.StatusBeforePause,
			&i.IsActive,
			&i.Type,
			&i.TotalPauseTime,
			&i.LastPing,
			&i.LastPausedAt,
			&i.LastResumedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserMonitors = `-- name: GetUserMonitors :many
SELECT id, name, period, grace_period, user_email, project_id, ping_url, status, status_before_pause, is_active, type, total_pause_time, last_ping, last_paused_at, last_resumed_at, created_at, updated_at FROM monitor WHERE user_email = $1
`

func (q *Queries) GetUserMonitors(ctx context.Context, userEmail string) ([]Monitor, error) {
	rows, err := q.db.Query(ctx, getUserMonitors, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Monitor
	for rows.Next() {
		var i Monitor
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Period,
			&i.GracePeriod,
			&i.UserEmail,
			&i.ProjectID,
			&i.PingUrl,
			&i.Status,
			&i.StatusBeforePause,
			&i.IsActive,
			&i.Type,
			&i.TotalPauseTime,
			&i.LastPing,
			&i.LastPausedAt,
			&i.LastResumedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserProjects = `-- name: GetUserProjects :many
SELECT id, name, visibility, user_email, created_at, updated_at FROM project WHERE user_email = $1
`

func (q *Queries) GetUserProjects(ctx context.Context, userEmail string) ([]Project, error) {
	rows, err := q.db.Query(ctx, getUserProjects, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Visibility,
			&i.UserEmail,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserUsingEmail = `-- name: GetUserUsingEmail :one
SELECT id, name, email, password, is_active, otp, last_login, created_at, updated_at FROM USERS WHERE email = $1
`

func (q *Queries) GetUserUsingEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserUsingEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.Otp,
		&i.LastLogin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserUsingOtp = `-- name: GetUserUsingOtp :one
SELECT id, name, email, password, is_active, otp, last_login, created_at, updated_at FROM USERS WHERE otp = $1
`

func (q *Queries) GetUserUsingOtp(ctx context.Context, otp *string) (User, error) {
	row := q.db.QueryRow(ctx, getUserUsingOtp, otp)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.Otp,
		&i.LastLogin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const pauseMonitor = `-- name: PauseMonitor :one
UPDATE monitor SET status = $1, status_before_pause = $2, last_paused_at = $3 WHERE id = $4 RETURNING id, name, period, grace_period, user_email, project_id, ping_url, status, status_before_pause, is_active, type, total_pause_time, last_ping, last_paused_at, last_resumed_at, created_at, updated_at
`

type PauseMonitorParams struct {
	Status            string
	StatusBeforePause *string
	LastPausedAt      pgtype.Timestamp
	ID                string
}

func (q *Queries) PauseMonitor(ctx context.Context, arg PauseMonitorParams) (Monitor, error) {
	row := q.db.QueryRow(ctx, pauseMonitor,
		arg.Status,
		arg.StatusBeforePause,
		arg.LastPausedAt,
		arg.ID,
	)
	var i Monitor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Period,
		&i.GracePeriod,
		&i.UserEmail,
		&i.ProjectID,
		&i.PingUrl,
		&i.Status,
		&i.StatusBeforePause,
		&i.IsActive,
		&i.Type,
		&i.TotalPauseTime,
		&i.LastPing,
		&i.LastPausedAt,
		&i.LastResumedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const resetUserPassword = `-- name: ResetUserPassword :exec
UPDATE USERS SET password = $1, otp = '' WHERE email = $2
`

type ResetUserPasswordParams struct {
	Password string
	Email    string
}

func (q *Queries) ResetUserPassword(ctx context.Context, arg ResetUserPasswordParams) error {
	_, err := q.db.Exec(ctx, resetUserPassword, arg.Password, arg.Email)
	return err
}

const resumeMonitor = `-- name: ResumeMonitor :one
UPDATE monitor m SET status = m.status_before_pause, status_before_pause = '', last_resumed_at = $1, total_pause_time = $2 WHERE id = $3 RETURNING id, name, period, grace_period, user_email, project_id, ping_url, status, status_before_pause, is_active, type, total_pause_time, last_ping, last_paused_at, last_resumed_at, created_at, updated_at
`

type ResumeMonitorParams struct {
	LastResumedAt  pgtype.Timestamp
	TotalPauseTime *int32
	ID             string
}

func (q *Queries) ResumeMonitor(ctx context.Context, arg ResumeMonitorParams) (Monitor, error) {
	row := q.db.QueryRow(ctx, resumeMonitor, arg.LastResumedAt, arg.TotalPauseTime, arg.ID)
	var i Monitor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Period,
		&i.GracePeriod,
		&i.UserEmail,
		&i.ProjectID,
		&i.PingUrl,
		&i.Status,
		&i.StatusBeforePause,
		&i.IsActive,
		&i.Type,
		&i.TotalPauseTime,
		&i.LastPing,
		&i.LastPausedAt,
		&i.LastResumedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateMonitor = `-- name: UpdateMonitor :exec
UPDATE monitor SET name = $1, period = $2, grace_period = $3 WHERE id = $4 AND user_email = $5
`

type UpdateMonitorParams struct {
	Name        string
	Period      int32
	GracePeriod int32
	ID          string
	UserEmail   string
}

func (q *Queries) UpdateMonitor(ctx context.Context, arg UpdateMonitorParams) error {
	_, err := q.db.Exec(ctx, updateMonitor,
		arg.Name,
		arg.Period,
		arg.GracePeriod,
		arg.ID,
		arg.UserEmail,
	)
	return err
}

const updateMonitorLastPing = `-- name: UpdateMonitorLastPing :exec
UPDATE monitor SET last_ping = $1, status='up' WHERE id = $2
`

type UpdateMonitorLastPingParams struct {
	LastPing pgtype.Timestamp
	ID       string
}

func (q *Queries) UpdateMonitorLastPing(ctx context.Context, arg UpdateMonitorLastPingParams) error {
	_, err := q.db.Exec(ctx, updateMonitorLastPing, arg.LastPing, arg.ID)
	return err
}

const updateMonitorStatus = `-- name: UpdateMonitorStatus :exec
UPDATE monitor SET status = $1 WHERE id = $2
`

type UpdateMonitorStatusParams struct {
	Status string
	ID     string
}

func (q *Queries) UpdateMonitorStatus(ctx context.Context, arg UpdateMonitorStatusParams) error {
	_, err := q.db.Exec(ctx, updateMonitorStatus, arg.Status, arg.ID)
	return err
}

const updateMonitorTotalPauseTime = `-- name: UpdateMonitorTotalPauseTime :exec
UPDATE monitor set total_pause_time = $1 where id = $2
`

type UpdateMonitorTotalPauseTimeParams struct {
	TotalPauseTime *int32
	ID             string
}

func (q *Queries) UpdateMonitorTotalPauseTime(ctx context.Context, arg UpdateMonitorTotalPauseTimeParams) error {
	_, err := q.db.Exec(ctx, updateMonitorTotalPauseTime, arg.TotalPauseTime, arg.ID)
	return err
}

const updateUserMonitorName = `-- name: UpdateUserMonitorName :exec
UPDATE monitor SET name = $1 WHERE user_email = $2
`

type UpdateUserMonitorNameParams struct {
	Name      string
	UserEmail string
}

func (q *Queries) UpdateUserMonitorName(ctx context.Context, arg UpdateUserMonitorNameParams) error {
	_, err := q.db.Exec(ctx, updateUserMonitorName, arg.Name, arg.UserEmail)
	return err
}

const updateUserMonitorSchedule = `-- name: UpdateUserMonitorSchedule :exec
UPDATE monitor SET period = $1, grace_period = $2 WHERE user_email = $3
`

type UpdateUserMonitorScheduleParams struct {
	Period      int32
	GracePeriod int32
	UserEmail   string
}

func (q *Queries) UpdateUserMonitorSchedule(ctx context.Context, arg UpdateUserMonitorScheduleParams) error {
	_, err := q.db.Exec(ctx, updateUserMonitorSchedule, arg.Period, arg.GracePeriod, arg.UserEmail)
	return err
}

const updateUserOtp = `-- name: UpdateUserOtp :exec
UPDATE USERS SET otp = $1 WHERE email = $2
`

type UpdateUserOtpParams struct {
	Otp   *string
	Email string
}

func (q *Queries) UpdateUserOtp(ctx context.Context, arg UpdateUserOtpParams) error {
	_, err := q.db.Exec(ctx, updateUserOtp, arg.Otp, arg.Email)
	return err
}

const updateUserProjectName = `-- name: UpdateUserProjectName :exec
UPDATE project SET name = $1 WHERE user_email = $2
`

type UpdateUserProjectNameParams struct {
	Name      string
	UserEmail string
}

func (q *Queries) UpdateUserProjectName(ctx context.Context, arg UpdateUserProjectNameParams) error {
	_, err := q.db.Exec(ctx, updateUserProjectName, arg.Name, arg.UserEmail)
	return err
}
