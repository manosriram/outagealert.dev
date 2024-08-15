package monitor

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
)

func Monitors(c echo.Context, env *types.Env) error {
	s, _ := session.Get("session", c)
	id := s.Values["id"].(int32)

	monitors, _ := env.DB.Query.GetUserMonitors(c.Request().Context(), id)
	fmt.Println(monitors)
	return c.Render(200, "monitors.html", template.UserMonitors{Monitors: monitors})
}

func Monitor(c echo.Context, env *types.Env) error {
	monitor_id := c.QueryParam("id")

	uu, _ := uuid.FromBytes([]byte(monitor_id))

	monitor, _ := env.DB.Query.GetMonitorById(c.Request().Context(), pgtype.UUID{Bytes: uu})
	return c.Render(200, "monitors.html", template.UserMonitor{Monitor: monitor})
}
