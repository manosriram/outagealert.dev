package ping

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

const (
	PING_HOST = "https://ping.outagealert.io"
)

func Ping(c echo.Context, env *types.Env) error {
	pingSlug := c.Param("ping_slug")
	url := fmt.Sprintf("%s/%s", PING_HOST, pingSlug)

	monitor, err := env.DB.Query.GetMonitorByPingUrl(c.Request().Context(), url)
	if err != nil {
		return c.JSON(500, "NOTOK")
	}

	err = env.DB.Query.CreatePing(c.Request().Context(), monitor.ID)
	if err != nil {
		return c.JSON(500, "NOTOK")
	}

	env.DB.Query.UpdateMonitorLastPing(c.Request().Context(), db.UpdateMonitorLastPingParams{LastPing: pgtype.Timestamp{Time: time.Now(), Valid: true}, ID: monitor.ID})

	return c.JSON(200, "OK")
}
