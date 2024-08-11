package main

import (
	"context"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manosriram/outagealert.io/pkg/auth"
	"github.com/manosriram/outagealert.io/pkg/dashboard"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := session.Get("session", c)
		if err != nil {
			c.Error(err)
		}

		email := s.Values["email"]
		if email == nil {
			return c.Redirect(302, "/signin")
		}
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.Renderer = template.NewTemplate()

	e.Static("/static", "static")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(middleware.Recover())
	// e.Use(middleware.Logger())

	conn, _ := pgx.Connect(context.TODO(), "user=postgres dbname=outagealertio sslmode=verify-full")
	dbconn := db.New(conn)
	env := types.NewEnv(dbconn)
	e.Use(types.InjectEnv(env))

	apiHandler := e.Group("/api")

	// Template handlers
	e.GET("/signin", auth.Signin)
	e.GET("/signup", auth.Signup)
	e.GET("/dashboard", dashboard.DashboardHome, IsAuthenticated)

	authApiHandler := apiHandler.Group("/auth")
	authApiHandler.POST("/signup", types.WithEnv(auth.SignUpApi))
	authApiHandler.POST("/signin", types.WithEnv(auth.SignInApi))

	e.Logger.Fatal(e.Start(":1323"))
}
