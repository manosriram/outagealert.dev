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

func ForgotPassword(c echo.Context) error {
	return c.Render(200, "forgot-password.html", nil)
}

func ConfirmOtp(c echo.Context) error {
	return c.Render(200, "confirm-otp.html", nil)
}

func ConfirmPassword(c echo.Context) error {
	return c.Render(200, "confirm-password.html", nil)
}
