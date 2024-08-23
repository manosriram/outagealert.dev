package ping

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	PING_HOST            = "https://ping.outagealert.io"
	NANOID_ALPHABET_LIST = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"
	NANOID_LENGTH        = 22
)

func StartMonitorCheck(monitor db.Monitor, env *types.Env) {
	fmt.Printf("ticker started for monitor %s; period: %d minute\n", monitor.ID, monitor.Period)
	// ticker := time.Tick(time.Minute * time.Duration(monitor.Period))
	ticker := time.Tick(time.Second * 10)

	for {
		select {
		case <-ticker:
			latestMonitor, err := env.DB.Query.GetMonitorById(context.Background(), monitor.ID)
			var status string
			oldStatus := latestMonitor.Status
			if oldStatus == "paused" {
				continue
			}
			if err != nil {
				log.Warnf("Error getting monitor by Id: %s", err.Error())
			}

			monitorUpDeadline := time.Now().Add(-time.Duration(time.Duration(latestMonitor.Period) * time.Minute)).UTC()
			monitorGraceDeadline := monitorUpDeadline.Add(-time.Duration(time.Duration(latestMonitor.GracePeriod) * time.Minute)).UTC()

			// fmt.Printf("name: %s -- utc createdat %s -- utc latestping %s\n", latestMonitor.Name, latestMonitor.CreatedAt.Time.UTC(), latestMonitor.LastPing.Time.UTC(),)
			fmt.Printf("monitorUpDeadline: %s -- monitorGraceDeadline: %s -- now %s\n", monitorUpDeadline, monitorGraceDeadline, time.Now().UTC())

			// Set monitor status to 'down' iff last_ping occurred before deadline OR monitor is created before deadline
			if (latestMonitor.LastPing.Time.UTC().Before(monitorUpDeadline) && latestMonitor.LastPing.Valid) || (!latestMonitor.LastPing.Valid && latestMonitor.CreatedAt.Time.UTC().Before(monitorUpDeadline)) {
				if (latestMonitor.LastPing.Time.UTC().Before(monitorGraceDeadline) && latestMonitor.LastPing.Valid) || (!latestMonitor.LastPing.Valid && latestMonitor.CreatedAt.Time.UTC().Before(monitorGraceDeadline)) {
					status = "down"
				} else {
					status = "grace_period"
				}

				// use where clause with email
				env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
					ID:     latestMonitor.ID,
					Status: status,
				})
				fmt.Printf("updating monitor status %s to %s\n", latestMonitor.Name, status)
			} else {
				status = "up"
				env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
					ID:     latestMonitor.ID,
					Status: status,
				})
				fmt.Printf("updating monitor status %s to up\n", latestMonitor.Name)
			}
			eventId, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
			if err != nil {
				log.Warnf("Error generating nanoid for event creation: %s\n", err.Error())
				continue
			}
			if status != oldStatus {
				err = env.DB.Query.CreateEvent(context.Background(), db.CreateEventParams{
					ID:         eventId,
					MonitorID:  latestMonitor.ID,
					FromStatus: oldStatus,
					ToStatus:   status,
				})
				if err != nil {
					log.Warnf("Error creating new event: %s\n", err.Error())
				}
			}
			fmt.Println("ticked ", monitor.ID)
		}
	}
}

func Ping(c echo.Context, env *types.Env) error {
	pingSlug := c.Param("ping_slug")
	url := fmt.Sprintf("%s/%s", PING_HOST, pingSlug)

	monitor, err := env.DB.Query.GetMonitorByPingUrl(c.Request().Context(), url)
	if err != nil {
		return c.JSON(500, "NOTOK")
	}

	id, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
	if err != nil {
		log.Warnf("Error creating nanoid: %s\n", err.Error())
		return c.JSON(500, "NOTOK")
	}
	err = env.DB.Query.CreatePing(c.Request().Context(), db.CreatePingParams{
		ID:        id,
		MonitorID: monitor.ID,
	})
	if err != nil {
		log.Warnf("Error creating ping: %s\n", err.Error())
		return c.JSON(500, "NOTOK")
	}

	err = env.DB.Query.UpdateMonitorLastPing(c.Request().Context(), db.UpdateMonitorLastPingParams{LastPing: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}, ID: monitor.ID})
	if err != nil {
		log.Warnf("Error updating monitor last ping: %s\n", err.Error())
	}
	eventId, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
	if err != nil {
		log.Warnf("Error creating ping: %s\n", err.Error())
		return c.JSON(500, "NOTOK")
	}
	if monitor.Status != "up" {
		err = env.DB.Query.CreateEvent(context.Background(), db.CreateEventParams{
			ID:         eventId,
			MonitorID:  monitor.ID,
			FromStatus: monitor.Status,
			ToStatus:   "up",
		})
		if err != nil {
			log.Warnf("Error creating new event: %s\n", err.Error())
		}
	}

	return c.JSON(200, "OK")
}
