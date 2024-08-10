package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func Login(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
