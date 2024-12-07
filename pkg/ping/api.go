package ping

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/event"
	"github.com/manosriram/outagealert.io/pkg/integration"
	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	PING_HOST            = "https://ping.outagealert.dev"
	NANOID_ALPHABET_LIST = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"
	NANOID_LENGTH        = 22
)

func CalculateMonitorStatus(dbMonitor *db.Monitor, env *types.Env) (string, error) {
	var status string
	oldStatus := dbMonitor.Status

	var period, gracePeriod time.Duration
	switch dbMonitor.PeriodText {
	case "minutes":
		period = time.Duration(dbMonitor.Period) * time.Minute
	case "hours":
		period = time.Duration(dbMonitor.Period) * time.Hour
	case "days":
		period = time.Duration(dbMonitor.Period) * (24 * time.Hour)
	}
	switch dbMonitor.GracePeriodText {
	case "minutes":
		gracePeriod = time.Duration(dbMonitor.GracePeriod) * time.Minute
	case "hours":
		gracePeriod = time.Duration(dbMonitor.GracePeriod) * time.Hour
	case "days":
		gracePeriod = time.Duration(dbMonitor.GracePeriod) * (24 * time.Hour)
	}

	monitorUpDeadline := time.Now().Add(-period).Add(-time.Duration(*dbMonitor.TotalPauseTime) * time.Second).UTC()
	monitorGraceDeadline := monitorUpDeadline.Add(-gracePeriod).UTC()

	fmt.Println(dbMonitor.ID, monitorUpDeadline, period)
	if oldStatus == "paused" || oldStatus == "down" {
		return oldStatus, nil
	}
	// Set monitor status to 'down' iff last_ping occurred before deadline OR monitor is created before deadline
	if (dbMonitor.LastPing.Time.UTC().Before(monitorUpDeadline) && dbMonitor.LastPing.Valid) || (!dbMonitor.LastPing.Valid && dbMonitor.CreatedAt.Time.UTC().Before(monitorUpDeadline)) {
		if (dbMonitor.LastPing.Time.UTC().Before(monitorGraceDeadline) && dbMonitor.LastPing.Valid) || (!dbMonitor.LastPing.Valid && dbMonitor.CreatedAt.Time.UTC().Before(monitorGraceDeadline)) {
			status = "down"
		} else {
			status = "grace_period"
		}

		p := int32(0)
		if status == "down" && oldStatus == "grace_period" || status == "up" && oldStatus == "down" {
			env.DB.Query.UpdateMonitorTotalPauseTime(context.Background(), db.UpdateMonitorTotalPauseTimeParams{
				ID:             dbMonitor.ID,
				TotalPauseTime: &p,
			})
		}

		// use where clause with email
		env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
			ID:     dbMonitor.ID,
			Status: status,
		})
	} else {
		status = "up"
		env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
			ID:     dbMonitor.ID,
			Status: status,
		})
	}
	fmt.Println(status, oldStatus)

	if status != oldStatus {
		err := event.CreateEvent(context.Background(), dbMonitor.ID, oldStatus, status, env)
		if err != nil {
			l.Log.Warnf("Error creating new event: %s\n", err.Error())
			return oldStatus, err
		}
	}
	return status, nil
}

func StartMonitorCheck(monitor db.Monitor, env *types.Env) {
	l.Log.Info("Started Monitor check for ", monitor.ID)
	ticker := time.NewTicker(time.Second * 10)
	done := make(chan struct{})

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dbMonitor, err := env.DB.Query.GetMonitorById(context.Background(), monitor.ID)
			if err != nil {
				l.Log.Errorf("Error calculating monitor status %s", err.Error())
			}
			status, err := CalculateMonitorStatus(&dbMonitor, env)

			// TODO: combine below 3 queries into 1
			emailIntegration, _ := env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
				MonitorID: dbMonitor.ID,
				AlertType: "email",
			})
			webhookIntegration, _ := env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
				MonitorID: dbMonitor.ID,
				AlertType: "webhook",
			})
			slackIntegration, _ := env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
				MonitorID: dbMonitor.ID,
				AlertType: "slack",
			})

			if status == "down" {
				if !emailIntegration.EmailAlertSent && emailIntegration.IsActive { // email alert enabled
					monitorLink := fmt.Sprintf("%s/monitor/%s/%s", os.Getenv("HOST_WITH_SCHEME"), dbMonitor.ProjectID, dbMonitor.ID)
					emailNotif := integration.EmailNotification{Email: dbMonitor.UserEmail, Env: *env, MonitorName: dbMonitor.Name, MonitorId: dbMonitor.ID, MonitorLink: monitorLink, NotificationType: integration.MONITOR_DOWN}
					emailNotif.SendAlert()
				}
				if !webhookIntegration.WebhookAlertSent && webhookIntegration.IsActive {
					webhookNotif := integration.WebhookNotification{Url: *webhookIntegration.AlertTarget, Env: *env, NotificationType: integration.MONITOR_DOWN, MonitorId: dbMonitor.ID, MonitorName: dbMonitor.Name}
					webhookNotif.SendAlert()
				}
				if !slackIntegration.SlackAlertSent && slackIntegration.IsActive {
					fmt.Println("sending down slack alert")
					monitorLink := fmt.Sprintf("%s/monitor/%s/%s", os.Getenv("HOST_WITH_SCHEME"), dbMonitor.ProjectID, dbMonitor.ID)
					slackNotif := integration.SlackNotification{Env: *env, NotificationType: integration.MONITOR_DOWN, MonitorName: dbMonitor.Name, MonitorId: dbMonitor.ID, UserEmail: dbMonitor.UserEmail, MonitorLink: monitorLink}
					slackNotif.SendAlert()
				}
			}
		case <-done:
			return
		}
	}
}

func Ping(c echo.Context, env *types.Env) error {
	pingSlug := c.Param("ping_slug")
	status := int32(200)
	metadata, err := json.Marshal(c.Request().Header.Get("User-Agent"))

	dbMonitor, err := env.DB.Query.GetMonitorByPingUrl(c.Request().Context(), pingSlug)
	if err != nil {
		l.Log.Error(err)
		status = 500
		err = env.DB.Query.CreatePing(c.Request().Context(), db.CreatePingParams{
			ID:        gonanoid.MustGenerate(NANOID_ALPHABET_LIST, NANOID_LENGTH),
			MonitorID: dbMonitor.ID,
			Status:    &status,
			Metadata:  metadata,
		})
		return c.JSON(500, "NOTOK")
	}

	if dbMonitor.Status == "paused" {
		status = 400
		err = env.DB.Query.CreatePing(c.Request().Context(), db.CreatePingParams{
			ID:        gonanoid.MustGenerate(NANOID_ALPHABET_LIST, NANOID_LENGTH),
			MonitorID: dbMonitor.ID,
			Status:    &status,
			Metadata:  metadata,
		})
		return c.JSON(400, "NOTOK:PAUSED_MONITOR")
	}

	webhookIntegration, _ := env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: dbMonitor.ID,
		AlertType: "webhook",
	})
	emailIntegration, _ := env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: dbMonitor.ID,
		AlertType: "email",
	})
	slackIntegration, _ := env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: dbMonitor.ID,
		AlertType: "slack",
	})

	id, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
	if err != nil {
		status = 500
		err = env.DB.Query.CreatePing(c.Request().Context(), db.CreatePingParams{
			ID:        id,
			MonitorID: dbMonitor.ID,
			Status:    &status,
			Metadata:  metadata,
		})
		return c.JSON(500, "NOTOK")
	}

	err = env.DB.Query.CreatePing(c.Request().Context(), db.CreatePingParams{
		ID:        id,
		MonitorID: dbMonitor.ID,
		Status:    &status,
		Metadata:  metadata,
	})
	if err != nil {
		l.Log.Warnf("Error creating ping: %s\n", err.Error())
		status = 500
		err = env.DB.Query.CreatePing(c.Request().Context(), db.CreatePingParams{
			ID:        id,
			MonitorID: dbMonitor.ID,
			Status:    &status,
			Metadata:  metadata,
		})
		return c.JSON(500, "NOTOK")
	}

	err = env.DB.Query.UpdateMonitorLastPing(c.Request().Context(), db.UpdateMonitorLastPingParams{LastPing: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}, ID: dbMonitor.ID})
	if err != nil {
		l.Log.Warnf("Error updating monitor last ping: %s\n", err.Error())
	}

	if dbMonitor.Status != "up" {
		err = event.CreateEvent(context.Background(), dbMonitor.ID, dbMonitor.Status, "up", env)
		if err != nil {
			l.Log.Warnf("Error creating new event: %s\n", err.Error())
		}
	}

	// Since the monitor will not be down after this ping, replicas cannot send duplicate emails
	if dbMonitor.Status == "down" {
		env.DB.Query.UpdateAlertSentFlag(context.Background(), db.UpdateAlertSentFlagParams{
			EmailAlertSent:   false,
			SlackAlertSent:   false,
			WebhookAlertSent: false,
			MonitorID:        dbMonitor.ID,
		})
		fmt.Println("isact = ", slackIntegration.IsActive)

		if emailIntegration.IsActive {
			monitorLink := fmt.Sprintf("%s/monitor/%s/%s", os.Getenv("HOST_WITH_SCHEME"), dbMonitor.ProjectID, dbMonitor.ID)
			emailNotif := integration.EmailNotification{Email: dbMonitor.UserEmail, Env: *env, MonitorName: dbMonitor.Name, MonitorLink: monitorLink, NotificationType: integration.MONITOR_UP}
			emailNotif.SendAlert()
		}
		if webhookIntegration.IsActive {
			webhookNotif := integration.WebhookNotification{Url: *webhookIntegration.AlertTarget, Env: *env, NotificationType: integration.MONITOR_UP}
			webhookNotif.SendAlert()
		}
		if slackIntegration.IsActive {
			monitorLink := fmt.Sprintf("%s/monitor/%s/%s", os.Getenv("HOST_WITH_SCHEME"), dbMonitor.ProjectID, dbMonitor.ID)
			slackNotif := integration.SlackNotification{Env: *env, NotificationType: integration.MONITOR_UP, MonitorName: dbMonitor.Name, MonitorId: dbMonitor.ID, UserEmail: dbMonitor.UserEmail, MonitorLink: monitorLink}
			slackNotif.SendAlert()
		}
	}

	return c.JSON(200, "OK")
}
