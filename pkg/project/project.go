package project

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
)

func Projects(c echo.Context, env *types.Env) error {
	s, _ := session.Get("session", c)
	email := s.Values["email"].(string)

	projects, _ := env.DB.Query.GetUserProjects(c.Request().Context(), email)
	fmt.Println(projects)
	return c.Render(200, "projects.html", template.UserProjects{Projects: projects})
}

func Monitor(c echo.Context, env *types.Env) error {
	monitor_id := c.QueryParam("id")

	uu, _ := uuid.FromBytes([]byte(monitor_id))

	monitor, _ := env.DB.Query.GetMonitorById(c.Request().Context(), pgtype.UUID{Bytes: uu})
	return c.Render(200, "monitors.html", template.UserMonitor{Monitor: monitor})
}
