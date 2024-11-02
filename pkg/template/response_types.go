package template

import (
	"time"

	"github.com/manosriram/outagealert.io/sqlc/db"
)

type ResponseMetadata struct {
	CreatedAtDistance string
	CurrentlyUpFor    int32
	CurrentlyDownFor  int32
	LastPing          float64
	IncidentsCount    int32
	UpDownTimeUnits   string
	LastPingTimeUnits string
}

type Response struct {
	Message  string
	Error    string
	Metadata ResponseMetadata
}

type MonitorMetadata struct {
	TotalPings                 int32
	TotalEvents                int32
	LastToStatusUpMonitorEvent time.Time
	LastPing                   time.Time
	MonitorCreated             time.Time
	EmailIntegration           bool
	SlackIntegration           bool
	WebhookIntegration         bool
	EmailIntegrationMetadata   db.AlertIntegration
	SlackIntegrationMetadata   db.AlertIntegration
	WebhookIntegrationMetadata db.AlertIntegration
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
	Monitors []db.Monitor
	Project  db.Project
}

type UserMonitorDetails struct {
	Response
	Monitor db.GetMonitorWithEventsByIdRow
}

type MonitorAlertIntegrations struct {
	EmailIntegrationEnabled   bool
	SlackIntegrationEnabled   bool
	WebhookIntegrationEnabled bool
	EmailIntegration          db.AlertIntegration
	SlackIntegration          db.AlertIntegration
	WebhookIntegration        db.AlertIntegration
}

type UserMonitor struct {
	Response
	Monitor                  db.Monitor
	MonitorMetadata          MonitorMetadata
	MonitorAlertIntegrations MonitorAlertIntegrations
}

type UserProjects struct {
	Response
	Projects []db.GetUserProjectsRow
}

type UserProject struct {
	Response
	Project db.Project
}

type MonitorEvents struct {
	Response
	Activity    []db.GetMonitorActivityPaginatedRow
	MonitorID   string
	CurrentPage int
	NextPage    int
	HasNextPage bool
}

type MonitorIntegrations struct {
	Integrations []db.AlertIntegration
}

type OrderCreatedResponse struct {
	PaymentSessionId string
	ENV              string
	Name             string
}
