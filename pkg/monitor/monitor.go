package monitor

import (
	"fmt"

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
	return c.Render(200, "monitor.html", template.UserMonitor{Monitor: monitor, Response: template.Response{Metadata: template.ResponseMetadata{CreatedAtDistance: createdAtDistance}}})
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
