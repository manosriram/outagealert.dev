package monitor

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

const (
	MAX_DISPLAY_TIME_IN_SECONDS int32 = 60
	MAX_DISPLAY_TIME_IN_MINUTES int32 = 60
	MAX_DISPLAY_TIME_IN_HOURS   int32 = 24
)

func ProjectMonitors(c echo.Context, env *types.Env) error {
	project_id := c.Param("project_id")
	monitors, err := env.DB.Query.GetProjectMonitors(c.Request().Context(), project_id)
	fmt.Println("e = ", err)
	return c.Render(200, "monitors.html", template.UserMonitors{Monitors: monitors, ProjectId: project_id})
}

func calculateRunningTime(monitor db.Monitor, response *template.Response, timeInSeconds, timeInMinutes, timeInHours float64) {
	switch monitor.Status {
	case "up", "grace_period":
		if timeInSeconds <= float64(MAX_DISPLAY_TIME_IN_SECONDS) {
			response.Metadata.CurrentlyUpFor = int32(timeInSeconds)
			if response.Metadata.CurrentlyUpFor > 1 {
				response.Metadata.UpDownTimeUnits = "seconds"
			} else {
				response.Metadata.UpDownTimeUnits = "second"
			}
		} else if timeInMinutes <= float64(MAX_DISPLAY_TIME_IN_MINUTES) {
			response.Metadata.CurrentlyUpFor = int32(timeInMinutes)
			if response.Metadata.CurrentlyUpFor > 1 {
				response.Metadata.UpDownTimeUnits = "minutes"
			} else {
				response.Metadata.UpDownTimeUnits = "minute"
			}
		} else {
			response.Metadata.CurrentlyUpFor = int32(timeInHours)
			if response.Metadata.CurrentlyUpFor > 1 {
				response.Metadata.UpDownTimeUnits = "hours"
			} else {
				response.Metadata.UpDownTimeUnits = "hour"
			}
		}
	case "down":
		if timeInSeconds <= float64(MAX_DISPLAY_TIME_IN_SECONDS) {
			response.Metadata.CurrentlyDownFor = int32(timeInSeconds)
			if response.Metadata.CurrentlyUpFor > 1 {
				response.Metadata.UpDownTimeUnits = "seconds"
			} else {
				response.Metadata.UpDownTimeUnits = "second"
			}
		} else if timeInMinutes <= float64(MAX_DISPLAY_TIME_IN_MINUTES) {
			response.Metadata.CurrentlyDownFor = int32(timeInMinutes)
			if response.Metadata.CurrentlyDownFor > 1 {
				response.Metadata.UpDownTimeUnits = "minutes"
			} else {
				response.Metadata.UpDownTimeUnits = "minute"
			}
		} else {
			response.Metadata.CurrentlyDownFor = int32(timeInHours)
			if response.Metadata.CurrentlyDownFor > 1 {
				response.Metadata.UpDownTimeUnits = "hours"
			} else {
				response.Metadata.UpDownTimeUnits = "hour"
			}
		}
	}
}

func Monitor(c echo.Context, env *types.Env) error {
	monitorId := c.Param("monitor_id")
	monitor, err := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	if err != nil {
	}
	createdAtDistance := FormatTimeAgo(monitor.CreatedAt.Time)
	var currentlyUpForTimeInSeconds, currentlyUpForTimeInMinutes, currentlyUpForTimeInHours float64
	var currentlyDownForTimeInSeconds, currentlyDownForTimeInMinutes, currentlyDownForTimeInHours float64
	event, _ := env.DB.Query.GetLastToStatusUpMonitorEvent(c.Request().Context(), monitorId)

	var status string
	switch monitor.Status {
	case "up", "grace_period":
		status = "up"
		currentlyUpForTimeInSeconds = time.Since(event.CreatedAt.Time).Seconds()
		currentlyUpForTimeInMinutes = time.Since(event.CreatedAt.Time).Minutes()
		currentlyUpForTimeInHours = time.Since(event.CreatedAt.Time).Hours()
	case "down":
		status = "down"
		down, _ := env.DB.Query.GetLatestMonitorEventByToStatus(c.Request().Context(), db.GetLatestMonitorEventByToStatusParams{
			MonitorID: monitorId,
			ToStatus:  "down",
		})
		currentlyDownForTimeInMinutes = math.Floor(time.Since(down.CreatedAt.Time).Minutes())
		currentlyDownForTimeInSeconds = math.Floor(time.Since(down.CreatedAt.Time).Seconds())
		currentlyDownForTimeInHours = math.Floor(time.Since(down.CreatedAt.Time).Hours())
	}

	response := template.Response{Metadata: template.ResponseMetadata{CreatedAtDistance: createdAtDistance, LastPing: -1, CurrentlyDownFor: -1, CurrentlyUpFor: -1}}
	var lastPingInSeconds, lastPingInMinutes, lastPingInHours float64 = -1, -1, -1
	if monitor.LastPing.Valid {
		lastPingInMinutes = math.Floor(time.Since(monitor.LastPing.Time).Minutes())
		lastPingInHours = math.Floor(time.Since(monitor.LastPing.Time).Hours())
		lastPingInSeconds = math.Floor(time.Since(monitor.LastPing.Time).Seconds())

		if lastPingInSeconds <= float64(MAX_DISPLAY_TIME_IN_SECONDS) {
			response.Metadata.LastPing = lastPingInSeconds
			response.Metadata.LastPingTimeUnits = "seconds"
		} else if lastPingInMinutes <= float64(MAX_DISPLAY_TIME_IN_MINUTES) {
			response.Metadata.LastPing = lastPingInMinutes
			response.Metadata.LastPingTimeUnits = "minutes"
		} else {
			response.Metadata.LastPing = lastPingInHours
			response.Metadata.LastPingTimeUnits = "hours"
		}
	}

	incidents, _ := env.DB.Query.GetNumberOfMonitorIncidents(c.Request().Context(), monitorId)
	response.Metadata.IncidentsCount = int32(incidents)

	totalPingCount, _ := env.DB.Query.TotalMonitorPings(c.Request().Context(), monitorId)
	totalEventCount, _ := env.DB.Query.TotalMonitorEvents(c.Request().Context(), monitorId)

	if status == "up" {
		calculateRunningTime(monitor, &response, currentlyUpForTimeInSeconds, currentlyUpForTimeInMinutes, currentlyUpForTimeInHours)
	} else {
		calculateRunningTime(monitor, &response, currentlyDownForTimeInSeconds, currentlyDownForTimeInMinutes, currentlyDownForTimeInHours)
	}
	pingUrl := url.URL{Host: os.Getenv("PING_HOST"), Scheme: os.Getenv("SCHEME"), Path: fmt.Sprintf("/%s", monitor.PingUrl)}

	// pingUrl := fmt.Sprintf("%s/%s", os.Getenv("PING_HOST"), monitor.PingUrl)
	monitor.PingUrl = pingUrl.String()
	fmt.Println(monitor.PingUrl)

	return c.Render(200, "monitor.html", template.UserMonitor{Monitor: monitor, Response: response, MonitorMetadata: template.MonitorMetadata{
		TotalPings:  int32(totalPingCount),
		TotalEvents: int32(totalEventCount),
	}})
}

// func MonitorPings(c echo.Context, env *types.Env) error {
// page := 1
// offset := 1
// monitorId := c.Param("monitor_id")
// events, err := env.DB.Query.GetPingsByMonitorIdPaginated(c.Request().Context(), db.GetPingsByMonitorIdPaginatedParams{
// MonitorID: monitorId,
// Offset:    int32(offset),
// })
// if err != nil {
// fmt.Println(err)
// }

// hasNextPage := true
// if len(events) == 0 {
// hasNextPage = false
// }
// return c.Render(200, "monitor-events-page.html", template.MonitorEvents{MonitorID: monitorId, Events: events, CurrentPage: page, NextPage: page + 1, HasNextPage: hasNextPage})
// }

func MonitorEvents(c echo.Context, env *types.Env) error {
	page := 1
	offset := 1
	monitorId := c.Param("monitor_id")
	activity, err := env.DB.Query.GetMonitorActivityPaginated(c.Request().Context(), db.GetMonitorActivityPaginatedParams{
		MonitorID: monitorId,
		Offset:    int32(offset),
	})
	if err != nil {
		fmt.Println(err)
	}

	hasNextPage := true
	if len(activity) == 0 {
		hasNextPage = false
	}
	return c.Render(200, "monitor-events-page.html", template.MonitorEvents{MonitorID: monitorId, Activity: activity, CurrentPage: page, NextPage: page + 1, HasNextPage: hasNextPage})
}
