package monitor

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
)

func ProjectMonitors(c echo.Context, env *types.Env) error {
	// s, _ := session.Get("session", c)
	// email := s.Values["email"].(string)

	project_id := c.Param("project_id")
	project_uuid, _ := uuid.Parse(project_id)
	fmt.Println(project_id, project_uuid)

	monitors, _ := env.DB.Query.GetProjectMonitors(c.Request().Context(), pgtype.UUID{Bytes: project_uuid, Valid: true})
	// for _, monitor := range monitors {
	// monitor.CreatedAt = monitor.CreatedAt.Time.Format("2006-01-02 15:04:05")
	// }
	return c.Render(200, "monitors.html", template.UserMonitors{Monitors: monitors, ProjectId: project_id})
}

func Monitor(c echo.Context, env *types.Env) error {
	monitor_id := c.QueryParam("id")

	uu, _ := uuid.FromBytes([]byte(monitor_id))

	monitor, _ := env.DB.Query.GetMonitorById(c.Request().Context(), pgtype.UUID{Bytes: uu})
	return c.Render(200, "monitors.html", template.UserMonitor{Monitor: monitor})
}
