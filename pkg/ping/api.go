package ping

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/manosriram/outagealert.io/pkg/event"
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
	ticker := time.Tick(time.Second * 1)

	for {
		select {
		case <-ticker:
			latestMonitor, err := env.DB.Query.GetMonitorById(context.Background(), monitor.ID)
			var status string
			oldStatus := latestMonitor.Status
			if err != nil {
				log.Warnf("Error getting monitor by Id: %s", err.Error())
			}

			/*
				get latest non paused event ie, to status != paused; and compare it with latest paused event ie, to_status = paused.
				check if non paused event occured after paused event. If yes, it means the
			*/

			// nonPausedEvent, err := env.DB.Query.GetLatestNonPausedMonitorEvent(context.Background(), latestMonitor.ID)
			// pausedEvent, err := env.DB.Query.GetLastToPausedMonitorEvent(context.Background(), latestMonitor.ID)
			// fmt.Println(latestMonitor.Name, upEvent.CreatedAt.Time, pausedEvent.CreatedAt.Time)

			// s1 := time.Since(upEvent.CreatedAt.Time.UTC()).Seconds()
			// s2 := time.Since(pausedEvent.CreatedAt.Time.UTC()).Seconds()

			// var z float64
			// // TODO: get latest pause session duration and add it to deadline
			// if latestMonitor.LastPausedAt.Valid == false || nonPausedEvent.CreatedAt.Time.After(pausedEvent.CreatedAt.Time) {
			// z = 0
			// } else {
			// z = nonPausedEvent.CreatedAt.Time.UTC().Sub(pausedEvent.CreatedAt.Time.UTC()).Seconds()
			// }

			// fmt.Println(s1, s2)
			// z := max(s1, s2) - min(s1, s2)
			// var dur float64
			// if latestMonitor.LastResumedAt.Valid && latestMonitor.LastPausedAt.Valid {
			// dur = latestMonitor.LastResumedAt.Time.Sub(latestMonitor.LastPausedAt.Time).Seconds()
			// } else {
			// dur = 0
			// }

			monitorUpDeadline := time.Now().Add(-time.Duration(latestMonitor.Period) * time.Minute).Add(-time.Duration(*latestMonitor.TotalPauseTime) * time.Second).UTC()
			monitorGraceDeadline := monitorUpDeadline.Add(-time.Duration(latestMonitor.GracePeriod) * time.Minute).UTC()

			fmt.Println("added paused seconds = ", *latestMonitor.TotalPauseTime)
			// fmt.Printf("%s\n%v -- %v\n%v\n\n", latestMonitor.Name, monitorUpDeadline, monitorGraceDeadline, latestMonitor.LastPing.Time.UTC())

			// monitorUpDeadline = monitorUpDeadline.Add(-time.Duration(z) * time.Second).UTC()
			// monitorGraceDeadline = monitorGraceDeadline.Add(-time.Duration(z) * time.Second).UTC()

			// fmt.Println(upEvent.CreatedAt.Time.UTC().Add(-time.Duration(pausedEvent.CreatedAt.Time.Second()) * time.Second).Second())
			// fmt.Printf("%s\n%v -- %v\n%v\n\n", latestMonitor.Name, monitorUpDeadline, monitorGraceDeadline, latestMonitor.LastPing.Time.UTC())

			// diff := upEvent.CreatedAt.Time.Sub(pausedEvent.CreatedAt.Time)

			// monitorUpDeadline = 2024-08-24 19:24:57.835851
			// monitorGraceDeadline = 2024-08-24 19:23:07.835565 +0000 UTC
			// paused = 60sec
			// last_paused_at = ...

			// fmt.Printf("%s\n%s\n%s\n\n", latestMonitor.Name, monitorUpDeadline, monitorGraceDeadline)
			// fmt.Printf("without added pause time deadline for %s is %s\n", latestMonitor.Name, latestMonitor.LastPing.Time.UTC())
			// fmt.Printf("added pause time deadline for %s is %s\n", latestMonitor.Name, latestMonitor.LastPing.Time.UTC().Add(-time.Duration(z)*time.Second))
			if oldStatus == "paused" {
				continue
			}

			// monitorUpDeadline := time.Now().Add(-time.Duration(time.Duration(latestMonitor.Period) * time.Minute)).Add(-time.Duration(*latestMonitor.TotalPauseTime) * time.Second).UTC()
			// monitorGraceDeadline := monitorUpDeadline.Add(-time.Duration(time.Duration(latestMonitor.GracePeriod) * time.Minute)).Add(-time.Duration(*latestMonitor.TotalPauseTime) * time.Second).UTC()

			// fmt.Printf("monitor %s\nperiod without pause time consideration: %v\nperiod with pause time consideration: %v\n\n", latestMonitor.Name, time.Now().Add(-time.Duration(time.Duration(latestMonitor.Period)*time.Minute)).UTC(), monitorUpDeadline)

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
				fmt.Printf("updating monitor status %s to %s\n", latestMonitor.Name, status)
			} else {
				status = "up"
				env.DB.Query.UpdateMonitorStatus(context.Background(), db.UpdateMonitorStatusParams{
					ID:     latestMonitor.ID,
					Status: status,
				})
				fmt.Printf("updating monitor status %s to up\n", latestMonitor.Name)
			}
			if status != oldStatus {
				err = event.CreateEvent(context.Background(), latestMonitor.ID, oldStatus, status, env)
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

	dbMonitor, err := env.DB.Query.GetMonitorByPingUrl(c.Request().Context(), url)
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
		MonitorID: dbMonitor.ID,
	})
	if err != nil {
		log.Warnf("Error creating ping: %s\n", err.Error())
		return c.JSON(500, "NOTOK")
	}

	err = env.DB.Query.UpdateMonitorLastPing(c.Request().Context(), db.UpdateMonitorLastPingParams{LastPing: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}, ID: dbMonitor.ID})
	if err != nil {
		log.Warnf("Error updating monitor last ping: %s\n", err.Error())
	}
	err = event.CreateEvent(context.Background(), dbMonitor.ID, dbMonitor.Status, "up", env)
	if err != nil {
		log.Warnf("Error creating new event: %s\n", err.Error())
	}

	return c.JSON(200, "OK")
}
