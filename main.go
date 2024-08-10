package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manosriram/outagealert.io/pkg/template"
)

func main() {
	e := echo.New()
	e.Renderer = template.NewTemplate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.Render(200, "home.html", nil)
}
