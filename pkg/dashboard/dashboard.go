package dashboard

import (
	"github.com/labstack/echo/v4"
)

func DashboardHome(c echo.Context) error {
	return c.Render(200, "dashboard.html", nil)
}

func Pricing(c echo.Context) error {
	return c.Render(200, "pricing.html", nil)
}

func Faq(c echo.Context) error {
	return c.Render(200, "faq.html", nil)
}

func Terms(c echo.Context) error {
	return c.Render(200, "termsconditions.html", nil)
}

func Contact(c echo.Context) error {
	return c.Render(200, "contact.html", nil)
}

func Refund(c echo.Context) error {
	return c.Render(200, "refund.html", nil)
}
