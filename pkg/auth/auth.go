package auth

import (
	"github.com/labstack/echo/v4"
)

func Signup(c echo.Context) error {
	return c.Render(200, "signup.html", nil)
}

func Signin(c echo.Context) error {
	return c.Render(200, "signin.html", nil)
}
