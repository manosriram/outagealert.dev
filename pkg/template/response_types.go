package template

import "github.com/manosriram/outagealert.io/sqlc/db"

type ResponseMetadata struct {
	CreatedAtDistance string
	CurrentlyUpFor    int32
	LastPing          float64
	IncidentsCount    int32
}

type Response struct {
	Message  string
	Error    string
	Metadata ResponseMetadata
}

type RegisterSuccessResponse struct {
	Response
	Email string
}

type ForgotPasswordSuccessResponse struct {
	Response
	Email string
}

type ResetPasswordResponse struct {
	Response
	Otp string
}

type User struct {
	Response
	User db.User
}

type UserMonitors struct {
	Response
	Monitors  []db.Monitor
	ProjectId string
}

type UserMonitorDetails struct {
	Response
	Monitor db.GetMonitorWithEventsByIdRow
}

type UserMonitor struct {
	Response
	Monitor db.Monitor
}

type UserProjects struct {
	Response
	Projects []db.Project
}

type UserProject struct {
	Response
	Project db.Project
}

type MonitorEvents struct {
	Response
	Events      []db.Event
	MonitorID   string
	CurrentPage int
	NextPage    int
	HasNextPage bool
}
