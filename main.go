package main

import (
	"context"
	"log"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manosriram/outagealert.io/pkg/auth"
	"github.com/manosriram/outagealert.io/pkg/monitor"
	"github.com/manosriram/outagealert.io/pkg/project"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

func ToDashboardIfAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := session.Get("session", c)
		if err != nil {
			c.Error(err)
		}

		email := s.Values["email"]
		if email == nil {
			return next(c)
		}
		return c.Redirect(302, "/monitors")
	}
}

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

	conn, err := pgx.Connect(context.TODO(), "user=postgres dbname=outagealertio sslmode=verify-full")
	if err != nil {
		log.Fatalf("Error connecting db: %s\n", err.Error())
	}
	dbconn := db.New(conn)
	env := types.NewEnv(dbconn)
	e.Use(types.InjectEnv(env))

	apiHandler := e.Group("/api")

	e.GET("/", func(c echo.Context) error {
		// return c.Redirect(302, "/signin")
		return c.Render(200, "signin.html", nil)
	})

	// Template handlers
	e.GET("/signin", auth.Signin, ToDashboardIfAuthenticated)
	e.GET("/signup", auth.Signup, ToDashboardIfAuthenticated)
	e.GET("/confirm-otp", auth.ConfirmOtp, ToDashboardIfAuthenticated)
	e.GET("/forgot-password", auth.ForgotPassword, ToDashboardIfAuthenticated)
	// e.GET("/confirm-password", auth.ConfirmPassword, ToDashboardIfAuthenticated)

	authApiHandler := apiHandler.Group("/auth")
	authApiHandler.POST("/signup", types.WithEnv(auth.SignUpApi))
	authApiHandler.POST("/signin", types.WithEnv(auth.SignInApi))
	authApiHandler.POST("/forgot-password", types.WithEnv(auth.ForgotPasswordApi))
	authApiHandler.POST("/confirm-otp", types.WithEnv(auth.ConfirmOtpApi))
	authApiHandler.POST("/reset-password", types.WithEnv(auth.ResetPasswordApi))

	// monitorApiHandler := apiHandler.Group("/monitors")
	// e.GET("/monitors", types.WithEnv(monitor.Monitors), IsAuthenticated)
	e.GET("/monitors/:project_id", types.WithEnv(monitor.ProjectMonitors), IsAuthenticated)
	e.POST("/api/monitors/create", types.WithEnv(monitor.CreateMonitor), IsAuthenticated)

	e.GET("/projects", types.WithEnv(project.Projects), IsAuthenticated)
	e.POST("/api/projects/create", types.WithEnv(project.CreateProject), IsAuthenticated)

	e.Logger.Fatal(e.Start(":1323"))
}
