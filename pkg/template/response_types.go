package template

import "github.com/manosriram/outagealert.io/sqlc/db"

type Response struct {
	Message string
	Error   string
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

type UserMonitors struct {
	Response
	Monitors  []db.Monitor
	ProjectId string
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
