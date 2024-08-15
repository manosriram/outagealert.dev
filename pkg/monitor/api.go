package monitor

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

type CreateMonitorForm struct {
	Name        string `form:"name" validate:"required"`
	Period      int32  `form:"period" validate:"required"`
	GracePeriod int32  `form:"grace_period" validate:"required"`
}

func CreateMonitor(c echo.Context, env *types.Env) error {
	createMonitorForm := new(CreateMonitorForm)
	if err := c.Bind(createMonitorForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	env.DB.Query.CreateMonitor(c.Request().Context(), db.CreateMonitorParams{Name: createMonitorForm.Name, Period: createMonitorForm.Period, GracePeriod: createMonitorForm.GracePeriod})
	fmt.Println("created")

	return nil
}
