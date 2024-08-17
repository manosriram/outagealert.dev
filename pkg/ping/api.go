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
	PING_HOST = "https://ping.outagealert.io"
)

func StartMonitorCheck(monitor db.GetAllMonitorIDsRow, env *types.Env) {
	fmt.Printf("ticker started for monitor %s; period: %d minute\n", monitor.ID, monitor.Period)
	// ticker := time.Tick(time.Second * time.Duration(monitor.Period))
	ticker := time.Tick(time.Second * 10)

	for {
		select {
		case <-ticker:
			latestMonitor, err := env.DB.Query.GetMonitorById(context.Background(), monitor.ID)
			if err != nil {
				log.Warnf("Error getting monitor by Id: %s", err.Error())
			}

			monitorUpDeadline := time.Now().Add(-time.Duration(time.Duration(latestMonitor.Period) * time.Minute)).UTC()
			var status string
			// if (latestMonitor.LastPing.Time.After(monitorUpDeadline) && latestMonitor.LastPing.Time.Before(time.Now())) || latestMonitor.CreatedAt.Time.Before(monitorUpDeadline) {

			fmt.Println(monitor.ID, latestMonitor.CreatedAt.Time.UTC(), monitorUpDeadline)

			// Set monitor status to 'down' iff last_ping occurred before deadline OR monitor is created before deadline
			if (latestMonitor.LastPing.Time.Before(monitorUpDeadline) && latestMonitor.LastPing.Valid) || latestMonitor.CreatedAt.Time.UTC().Before(monitorUpDeadline) {
				fmt.Println("to down")
				status = "down"
				env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
					ID:     latestMonitor.ID,
					Status: &status,
				})
			} else {
				fmt.Println("to up")
				status = "up"
				env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
					ID:     latestMonitor.ID,
					Status: &status,
				})
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

	id, err := gonanoid.New()
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

	return c.JSON(200, "OK")
}
