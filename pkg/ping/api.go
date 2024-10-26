package ping

import (
	"context"
	"encoding/json"
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

func StartMonitorCheck(monitor db.Monitor, env *types.Env) {
	ticker := time.NewTicker(time.Second * 10)
	done := make(chan struct{})

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			latestMonitor, err := env.DB.Query.GetMonitorById(context.Background(), monitor.ID)
			var status string
			oldStatus := latestMonitor.Status
			if err != nil {
				l.Log.Errorf("Ticker %s exiting", latestMonitor.ID)
				close(done)
				return
			}

			monitorUpDeadline := time.Now().Add(-time.Duration(latestMonitor.Period) * time.Minute).Add(-time.Duration(*latestMonitor.TotalPauseTime) * time.Second).UTC()
			monitorGraceDeadline := monitorUpDeadline.Add(-time.Duration(latestMonitor.GracePeriod) * time.Minute).UTC()

			if oldStatus == "paused" {
				continue
			}
			// Set monitor status to 'down' iff last_ping occurred before deadline OR monitor is created before deadline
			if (latestMonitor.LastPing.Time.UTC().Before(monitorUpDeadline) && latestMonitor.LastPing.Valid) || (!latestMonitor.LastPing.Valid && latestMonitor.CreatedAt.Time.UTC().Before(monitorUpDeadline)) {
				if (latestMonitor.LastPing.Time.UTC().Before(monitorGraceDeadline) && latestMonitor.LastPing.Valid) || (!latestMonitor.LastPing.Valid && latestMonitor.CreatedAt.Time.UTC().Before(monitorGraceDeadline)) {
					status = "down"
				} else {
					status = "grace_period"
				}

				/*
					up -> down
					down -> up
					down -> paused
					down -> grace_period
					grace_period -> down
					grace_period -> paused
					paused -> up
					paused -> down
					paused -> grace_period
				*/
				p := int32(0)
				if status == "down" && oldStatus == "grace_period" || status == "up" && oldStatus == "down" {
					env.DB.Query.UpdateMonitorTotalPauseTime(context.Background(), db.UpdateMonitorTotalPauseTimeParams{
						ID:             latestMonitor.ID,
						TotalPauseTime: &p,
					})
				}

				// use where clause with email
				env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
					ID:     latestMonitor.ID,
					Status: status,
				})
			} else {
				status = "up"
				env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
					ID:     latestMonitor.ID,
					Status: status,
				})
			}
			emailIntegration, err := env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
				MonitorID: latestMonitor.ID,
				AlertType: "email",
			})
			webhookIntegration, err := env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
				MonitorID: latestMonitor.ID,
				AlertType: "webhook",
			})

			if status != oldStatus {
				err = event.CreateEvent(context.Background(), latestMonitor.ID, oldStatus, status, env)
				if err != nil {
					l.Log.Warnf("Error creating new event: %s\n", err.Error())
				}
			}

			if status == "down" {
				if !emailIntegration.EmailAlertSent && emailIntegration.IsActive { // email alert enabled
					emailNotif := integration.EmailNotification{Email: latestMonitor.UserEmail, Env: *env}
					emailNotif.SendAlert(latestMonitor.ID, latestMonitor.Name)
				}
				if !webhookIntegration.WebhookAlertSent && webhookIntegration.IsActive {
					webhookNotif := integration.WebhookNotification{Url: *webhookIntegration.AlertTarget, Env: *env}
					webhookNotif.SendAlert(latestMonitor.ID, latestMonitor.Name)
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

	if dbMonitor.Status == "down" {
		env.DB.Query.UpdateAlertSentFlag(context.Background(), db.UpdateAlertSentFlagParams{
			EmailAlertSent:   false,
			SlackAlertSent:   false,
			WebhookAlertSent: false,
			MonitorID:        dbMonitor.ID,
		})
	}

	return c.JSON(200, "OK")
}
