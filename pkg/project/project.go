package project

import (
	"fmt"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/monitor"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
)

func Projects(c echo.Context, env *types.Env) error {
	s, _ := session.Get("session", c)
	email := s.Values["email"]
	if email == nil {
		host := os.Getenv("HOST_WITH_SCHEME")
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s", host))
		return c.Render(200, "projects.html", template.UserProjects{Response: template.Response{Error: "Access denied"}})
	}

	projects, err := env.DB.Query.GetUserProjects(c.Request().Context(), email.(string))
	if err != nil {
		l.Log.Errorf("Error getting user projects %s", err.Error())
		return c.Render(200, "projects.html", template.UserProjects{Response: template.Response{Error: "Internal server error"}})
	}

	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), email.(string))
	if err != nil {
		l.Log.Errorf("Error getting user %s", err.Error())
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.UserProjects{Response: template.Response{Error: "Error getting user"}})
	}

	monitorCount, err := env.DB.Query.GetTotalMonitorCount(c.Request().Context(), email.(string))
	if err != nil {
		l.Log.Error("Error getting monitor count ", err.Error())
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.UserProjects{Response: template.Response{Error: "Error getting monitor count"}})
	}

	var monitorLimit int64 = monitor.PLAN_VS_MONITOR_COUNT[*user.Plan]
	return c.Render(200, "projects.html", template.UserProjects{Projects: projects, MonitorLimit: monitorLimit, MonitorUsed: monitorCount[0]})
}

func Monitor(c echo.Context, env *types.Env) error {
	monitorId := c.QueryParam("id")

	monitor, err := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	if err != nil {
		host := os.Getenv("HOST_WITH_SCHEME")
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s", host))
		return c.Render(200, "monitors.html", template.UserMonitors{Response: template.Response{Error: "Internal server error"}})
	}
	return c.Render(
		200,
		"monitors.html",
		template.UserMonitor{Monitor: monitor},
	)
}
