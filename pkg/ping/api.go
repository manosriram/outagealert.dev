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

func StartMonitorCheck(monitor_id string, period, gracePeriod int32) {
	fmt.Printf("ticker started for monitor %s; period: %d minute\n", monitor_id, period)
	ticker := time.Tick(time.Second * time.Duration(period))

	for {
		select {
		case <-ticker:
			fmt.Println("ticked ", monitor_id)
		}
	}
}

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
