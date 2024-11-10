package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manosriram/outagealert.io/pkg/auth"
	"github.com/manosriram/outagealert.io/pkg/dashboard"
	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/monitor"
	"github.com/manosriram/outagealert.io/pkg/payment"
	"github.com/manosriram/outagealert.io/pkg/ping"
	"github.com/manosriram/outagealert.io/pkg/project"
	t "github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

func NotifyOutageAlert() {
	l.Log.Infof("Starting notifier")
	ticker := time.NewTicker(30 * time.Minute)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			urls := strings.Split(os.Getenv("MONITORING_URLS"), ";")
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					l.Log.Errorf("Error requesting monitoring url %s", err.Error())
					return
				}
				defer resp.Body.Close()

				_, err = io.ReadAll(resp.Body)
				if err != nil {
					l.Log.Errorf("Error requesting monitoring url %s", err.Error())
					return
				}

				l.Log.Infof("Notified %s", url)
			}

		case <-quit:
			ticker.Stop()
			return
		}
	}
}

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
		return c.Redirect(302, "/projects")
	}
}

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := session.Get("session", c)
		if err != nil {
			c.Error(err)
			return c.Redirect(302, "/signin")
		}

		email := s.Values["email"]
		if email == nil {
			// return c.Render(200, "signin.html", nil)
			return c.Redirect(302, "/signin")
		}
		return next(c)
	}
}

func initDB() *db.Queries {
	psqlUser := os.Getenv("POSTGRES_USER")
	psqlPassword := os.Getenv("POSTGRES_PASSWORD")
	psqlPort := os.Getenv("POSTGRES_PORT")
	psqlDatabase := os.Getenv("POSTGRES_DATABASE")
	psqlHost := os.Getenv("POSTGRES_HOST")

	psqlString := fmt.Sprintf("user=%s password=%s port=%s database=%s sslmode=disable host=%s", psqlUser, psqlPassword, psqlPort, psqlDatabase, psqlHost)
	l.Log.Info("psql string ", psqlString)
	config, err := pgxpool.ParseConfig(psqlString)
	if err != nil {
		l.Log.Infof("Unable to parse connection string: %v", err)
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	dbconn := db.New(pool)
	return dbconn
}

func main() {
	e := echo.New()
	e.Renderer = t.NewTemplate()
	l.Init() // Initialize logger

	rateLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      5,
				Burst:     30,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	e.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e.Static("/static", "static")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	dbconn := initDB()
	env := types.NewEnv(dbconn)
	e.Use(types.InjectEnv(env))

	apiHandler := e.Group("/api")

	// Start checks for monitors
	go monitor.StartAllMonitorChecks(env)

	// Template handlers
	e.GET("/", auth.Signin, ToDashboardIfAuthenticated)
	e.GET("/signin", auth.Signin, ToDashboardIfAuthenticated)
	e.GET("/signup", auth.Signup, ToDashboardIfAuthenticated)
	e.GET("/confirm-otp", auth.ConfirmOtp, ToDashboardIfAuthenticated)
	e.GET("/forgot-password", auth.ForgotPassword, ToDashboardIfAuthenticated)
	e.GET("/pricing", dashboard.Pricing)
	e.GET("/faq", dashboard.Faq)
	e.GET("/terms", dashboard.Terms)
	e.GET("/contact", dashboard.Contact)
	e.GET("/refund", dashboard.Refund)
	e.GET("/email-verified", dashboard.EmailVerified)

	authApiHandler := apiHandler.Group("/auth")
	authApiHandler.POST("/signup", types.WithEnv(auth.SignUpApi))
	authApiHandler.POST("/signin", types.WithEnv(auth.SignInApi))
	authApiHandler.POST("/forgot-password", types.WithEnv(auth.ForgotPasswordApi))
	authApiHandler.POST("/confirm-otp", types.WithEnv(auth.ConfirmOtpApi))
	authApiHandler.POST("/reset-password", types.WithEnv(auth.ResetPasswordApi))
	e.GET("/verify/:magic_token", types.WithEnv(auth.VerifyEmailViaMagicToken))

	e.GET("/user", types.WithEnv(auth.GetCurrentUser))
	e.GET("/user/logout", types.WithEnv(auth.Logout), IsAuthenticated)
	e.GET("/monitors/:project_id", types.WithEnv(monitor.ProjectMonitors), IsAuthenticated)
	e.GET("/monitor/:project_id/:monitor_id", types.WithEnv(monitor.Monitor), IsAuthenticated)
	e.GET("/monitor/:monitor_id/events", types.WithEnv(monitor.MonitorEvents), IsAuthenticated)
	e.GET("/monitor/:monitor_id/integrations", types.WithEnv(monitor.MonitorIntegrations), IsAuthenticated)
	e.PUT("/monitor/:monitor_id/integrations", types.WithEnv(monitor.UpdateMonitorIntegrations), IsAuthenticated)

	monitorApiHandler := apiHandler.Group("/monitor")
	monitorApiHandler.GET("/:monitor_id/table/events", types.WithEnv(monitor.GetMonitorEventsTable), IsAuthenticated)
	monitorApiHandler.GET("/:monitor_id/events", types.WithEnv(monitor.GetMonitorActivity), IsAuthenticated)
	monitorApiHandler.GET("/pause/:monitor_id", types.WithEnv(monitor.PauseMonitor), IsAuthenticated)
	monitorApiHandler.GET("/resume/:monitor_id", types.WithEnv(monitor.ResumeMonitor), IsAuthenticated)
	monitorApiHandler.PUT("/:monitor_id", types.WithEnv(monitor.UpdateMonitor), IsAuthenticated)
	monitorApiHandler.DELETE("/:project_id/:monitor_id", types.WithEnv(monitor.DeleteMonitor), IsAuthenticated)

	e.POST("/api/monitors/create", types.WithEnv(monitor.CreateMonitor), IsAuthenticated)
	e.POST("/api/projects/create", types.WithEnv(project.CreateProject), IsAuthenticated)
	e.GET("/projects", types.WithEnv(project.Projects), IsAuthenticated)

	e.GET("/p/:ping_slug", types.WithEnv(ping.Ping))
	e.POST("/api/contactus", types.WithEnv(dashboard.SubmitContact))

	// Payment APIs
	e.GET("/payment/create_order", types.WithEnv(payment.CreateOrder))
	e.POST("/payment-webhook", types.WithEnv(payment.OrderWebhook))

	/*
		Start the monitoring service
		Pings the outagealert monitor to notify liveness every 20 minutes
	*/
	// go NotifyOutageAlert()

	l.Log.Info("Starting server at :1323")
	e.Logger.Fatal(e.Start(":1323"))
}
