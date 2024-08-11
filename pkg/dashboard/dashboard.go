package dashboard

import (
	"github.com/labstack/echo/v4"
)

func DashboardHome(c echo.Context) error {
	return c.Render(200, "dashboard.html", nil)
}
