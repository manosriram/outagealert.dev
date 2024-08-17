package monitor

import (
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
)

func ProjectMonitors(c echo.Context, env *types.Env) error {
	project_id := c.Param("project_id")
	monitors, _ := env.DB.Query.GetProjectMonitors(c.Request().Context(), project_id)
	return c.Render(200, "monitors.html", template.UserMonitors{Monitors: monitors, ProjectId: project_id})
}

func Monitor(c echo.Context, env *types.Env) error {
	monitorId := c.Param("monitor_id")
	monitor, _ := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	return c.Render(200, "monitor.html", template.UserMonitor{Monitor: monitor})
}
