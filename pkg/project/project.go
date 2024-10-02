package project

import (
	"fmt"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
)

func Projects(c echo.Context, env *types.Env) error {
	s, _ := session.Get("session", c)
	email := s.Values["email"]
	if email == nil {
		// return c.Render(200, "errors", template.Response{Error: "Access denied"})
		host := os.Getenv("HOST_WITH_SCHEME")
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s", host))
		return c.Render(200, "projects.html", template.UserProjects{Response: template.Response{Error: "Access denied"}})
	}

	projects, err := env.DB.Query.GetUserProjects(c.Request().Context(), email.(string))
	if err != nil {
		return c.Render(200, "projects.html", template.UserProjects{Response: template.Response{Error: "Internal server error"}})
	}
	return c.Render(200, "projects.html", template.UserProjects{Projects: projects})
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
