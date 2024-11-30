package monitor

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/l"
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
	if err != nil {
		l.Log.Error("Error getting monitors ", err.Error())
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.UserProjects{Response: template.Response{Error: "Error getting monitors"}})
	}
	project, err := env.DB.Query.GetProjectById(context.Background(), project_id)
	if err != nil {
		l.Log.Error("Error getting projects ", err.Error())
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.UserProjects{Response: template.Response{Error: "Error getting projects"}})
	}

	return c.Render(200, "monitors.html", template.UserMonitors{Monitors: monitors, Project: project})
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
	host := os.Getenv("HOST_WITH_SCHEME")
	monitorId := c.Param("monitor_id")
	s, _ := session.Get("session", c)
	email := s.Values["email"]

	monitor, err := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/signin", host))
		return c.NoContent(400)
	}

	event, err := env.DB.Query.GetLastToStatusUpMonitorEvent(c.Request().Context(), monitorId)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/signin", host))
		return c.NoContent(400)
	}

	response := template.Response{Metadata: template.ResponseMetadata{}}

	incidents, err := env.DB.Query.GetNumberOfMonitorIncidents(c.Request().Context(), monitorId)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/signin", host))
		return c.NoContent(400)
	}
	response.Metadata.IncidentsCount = int32(incidents)

	totalPingCount, err := env.DB.Query.TotalMonitorPings(c.Request().Context(), monitorId)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/signin", host))
		return c.NoContent(400)
	}

	totalEventCount, err := env.DB.Query.TotalMonitorEvents(c.Request().Context(), monitorId)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/signin", host))
		return c.NoContent(400)
	}

	pingUrl := url.URL{Host: os.Getenv("PING_HOST"), Scheme: os.Getenv("SCHEME"), Path: fmt.Sprintf("/p/%s", monitor.PingUrl)}
	monitor.PingUrl = pingUrl.String()

	monitorAlertIntegrations := template.MonitorAlertIntegrations{}
	integrations, err := env.DB.Query.GetMonitorIntegrations(c.Request().Context(), monitorId)
	for _, integration := range integrations {
		switch integration.AlertType {
		case "email":
			monitorAlertIntegrations.EmailIntegrationEnabled = integration.IsActive
			monitorAlertIntegrations.EmailIntegration = integration
		case "slack":
			monitorAlertIntegrations.SlackIntegrationEnabled = integration.IsActive
			monitorAlertIntegrations.SlackIntegration = integration
			monitorAlertIntegrations.SlackAuthUrl = fmt.Sprintf(os.Getenv("SLACK_OAUTH_URL"), email)
		case "webhook":
			monitorAlertIntegrations.WebhookIntegrationEnabled = integration.IsActive
			monitorAlertIntegrations.WebhookIntegration = integration
		}
	}

	return c.Render(200, "monitor.html", template.UserMonitor{Monitor: monitor, MonitorAlertIntegrations: monitorAlertIntegrations, Response: response, MonitorMetadata: template.MonitorMetadata{
		MonitorCreated:             monitor.CreatedAt.Time,
		TotalPings:                 int32(totalPingCount),
		TotalEvents:                int32(totalEventCount),
		LastToStatusUpMonitorEvent: event.CreatedAt.Time,
		LastPing:                   monitor.LastPing.Time,
	}})
}

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
