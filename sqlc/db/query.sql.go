// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
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

const createMonitor = `-- name: CreateMonitor :one
INSERT INTO monitor(id, name, period, grace_period, user_email, project_id, ping_url, type) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, name, period, grace_period, user_email, project_id, ping_url, status, type, last_ping, created_at, updated_at
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
		&i.Type,
		&i.LastPing,
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

const getMonitorById = `-- name: GetMonitorById :one
SELECT id, name, period, grace_period, user_email, project_id, ping_url, status, type, last_ping, created_at, updated_at FROM monitor WHERE id = $1
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
		&i.Type,
		&i.LastPing,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMonitorByPingUrl = `-- name: GetMonitorByPingUrl :one
SELECT id, name, period, grace_period, user_email, project_id, ping_url, status, type, last_ping, created_at, updated_at FROM monitor WHERE ping_url = $1
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
		&i.Type,
		&i.LastPing,
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
SELECT id, name, period, grace_period, user_email, project_id, ping_url, status, type, last_ping, created_at, updated_at FROM monitor WHERE project_id = $1
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
			&i.Type,
			&i.LastPing,
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
SELECT id, name, period, grace_period, user_email, project_id, ping_url, status, type, last_ping, created_at, updated_at FROM monitor WHERE user_email = $1
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
			&i.Type,
			&i.LastPing,
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
