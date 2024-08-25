package monitor

import (
	"fmt"
	"math"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

func ProjectMonitors(c echo.Context, env *types.Env) error {
	project_id := c.Param("project_id")
	monitors, err := env.DB.Query.GetProjectMonitors(c.Request().Context(), project_id)
	fmt.Println("e = ", err)
	return c.Render(200, "monitors.html", template.UserMonitors{Monitors: monitors, ProjectId: project_id})
}

func Monitor(c echo.Context, env *types.Env) error {
	monitorId := c.Param("monitor_id")
	monitor, _ := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	createdAtDistance := formatTimeAgo(monitor.CreatedAt.Time)
	fmt.Println(createdAtDistance)
	var currentlyUpForTime float64
	event, _ := env.DB.Query.GetLastToStatusUpMonitorEvent(c.Request().Context(), monitorId)

	switch monitor.Status {
	case "up", "grace_period":
		currentlyUpForTime = float64(time.Now().Add(-time.Duration(event.CreatedAt.Time.UTC().Minute()) * time.Minute).Minute())
	case "down":
		down, _ := env.DB.Query.GetLatestMonitorEventByToStatus(c.Request().Context(), db.GetLatestMonitorEventByToStatusParams{
			MonitorID: monitorId,
			ToStatus:  "down",
		})
		currentlyUpForTime = float64(down.CreatedAt.Time.Add(-time.Duration(event.CreatedAt.Time.UTC().Minute()) * time.Minute).Minute())
	}

	lastPing := math.Floor(time.Since(monitor.LastPing.Time).Minutes())
	incidents, _ := env.DB.Query.GetNumberOfMonitorIncidents(c.Request().Context(), monitorId)
	return c.Render(200, "monitor.html", template.UserMonitor{Monitor: monitor, Response: template.Response{Metadata: template.ResponseMetadata{CreatedAtDistance: createdAtDistance, LastPing: lastPing, CurrentlyUpFor: int32(currentlyUpForTime), IncidentsCount: int32(incidents)}}})
}

func MonitorEvents(c echo.Context, env *types.Env) error {
	page := 1
	offset := 1
	monitorId := c.Param("monitor_id")
	events, err := env.DB.Query.GetEventsByMonitorIdPaginated(c.Request().Context(), db.GetEventsByMonitorIdPaginatedParams{
		MonitorID: monitorId,
		Offset:    int32(offset),
	})
	if err != nil {
		fmt.Println(err)
	}
	// for _, event := range events {

	// }

	hasNextPage := true
	if len(events) == 0 {
		hasNextPage = false
	}
	return c.Render(200, "monitor-events-page.html", template.MonitorEvents{MonitorID: monitorId, Events: events, CurrentPage: page, NextPage: page + 1, HasNextPage: hasNextPage})
}
