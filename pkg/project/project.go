package project

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
)

func Projects(c echo.Context, env *types.Env) error {
	s, _ := session.Get("session", c)
	email := s.Values["email"].(string)

	projects, err := env.DB.Query.GetUserProjects(c.Request().Context(), email)
	if err != nil {
		return c.Render(200, "error", template.Response{Error: ""})
	}
	return c.Render(200, "projects.html", template.UserProjects{Projects: projects})
}

func Monitor(c echo.Context, env *types.Env) error {
	monitorId := c.QueryParam("id")

	monitor, err := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	if err != nil {
		return c.Render(200, "error", template.Response{Error: ""})
	}
	return c.Render(200, "monitors.html", template.UserMonitorDetails{Monitor: monitor})
}
