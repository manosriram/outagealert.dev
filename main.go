package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	"github.com/manosriram/outagealert.io/pkg/ping"
	"github.com/manosriram/outagealert.io/pkg/project"
	"github.com/manosriram/outagealert.io/pkg/template"
	t "github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/pkg/webhook"
	"github.com/manosriram/outagealert.io/sqlc/db"
	"github.com/plutov/paypal"
	"golang.org/x/time/rate"
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
		return c.Redirect(302, "/projects")
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
			// return c.Render(200, "signin.html", nil)
			return c.Redirect(302, "/signin")
		}
		return next(c)
	}
}

var (
	clientID     = "AfXmSg_EyAbi0J13tuX5yy-ErT59xqXntYmRkSRF74JluM1Fnu7Q-1eOdBvVqKueD7cVWvOaKUxvHIc5"
	clientSecret = "EKvyV6-urgG7lZZukrm9mUu_lwvAT2K-DGgAPRcva1m_jNB-8uayFOuLVUxqetRxWYEzXxonX60phRp-"
	paypalClient *paypal.Client
)

// func handleIndex(e echo.Context) {
// tmpl, err := template.ParseFiles("index.html")
// if err != nil {
// // http.Error(w, err.Error(), http.StatusInternalServerError)
// return
// }
// tmpl.Execute(w, nil)
// }

type X struct {
	Amount string `form:"amount"`
}

func handlePayment(c echo.Context) error {
	// if r.Method != http.MethodPost {
	// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// return
	// }

	// amount := c.FormValue("amount")
	x := new(X)
	if err := c.Bind(x); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(x.Amount)

	// Create PayPal order
	order, _ := paypalClient.CreateOrder(
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Value:    x.Amount,
					Currency: "USD",
				},
			},
		},
		nil,
		nil,
	)
	fmt.Println(order.ID)

	// if err != nil {
	// http.Error(w, err.Error(), http.StatusInternalServerError)
	// return
	// }

	// Redirect to PayPal checkout
	for _, link := range order.Links {
		if link.Rel == "approve" {
			c.Response().Header().Set("HX-Redirect", link.Href)
			return nil
		}
	}

	return nil
	// http.Error(w, "Failed to create PayPal order", http.StatusInternalServerError)
}

func main() {
	e := echo.New()
	e.Renderer = t.NewTemplate()
	l.Init()

	rateLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.Render(http.StatusTooManyRequests, "errors", template.Response{Error: "Too many requests"})
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
	// e.Use(middleware.Logger())

	psqlUser := os.Getenv("POSTGRES_USER")
	psqlPassword := os.Getenv("POSTGRES_PASSWORD")
	psqlPort := os.Getenv("POSTGRES_PORT")
	psqlDatabase := os.Getenv("POSTGRES_DATABASE")
	psqlHost := os.Getenv("POSTGRES_HOST")

	config, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s port=%s database=%s sslmode=disable host=%s", psqlUser, psqlPassword, psqlPort, psqlDatabase, psqlHost))
	if err != nil {
		l.Log.Infof("Unable to parse connection string: %v", err)
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	dbconn := db.New(pool)
	env := types.NewEnv(dbconn)
	e.Use(types.InjectEnv(env))

	apiHandler := e.Group("/api")

	go monitor.StartAllMonitorChecks(env)

	// Template handlers
	e.GET("/", auth.Signin, ToDashboardIfAuthenticated) // TODO: redirect this to landing page
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

	e.POST("/webhook/paypal", webhook.PaypalWebhook)

	e.GET("/user", types.WithEnv(auth.GetCurrentUser))
	e.GET("/user/logout", types.WithEnv(auth.Logout), IsAuthenticated)
	e.GET("/monitors/:project_id", types.WithEnv(monitor.ProjectMonitors), IsAuthenticated)
	e.GET("/monitor/:project_id/:monitor_id", types.WithEnv(monitor.Monitor), IsAuthenticated)
	e.GET("/api/monitor/:monitor_id/events", types.WithEnv(monitor.GetMonitorActivity), IsAuthenticated)
	e.GET("/api/monitor/:monitor_id/table/events", types.WithEnv(monitor.GetMonitorEventsTable), IsAuthenticated)
	e.GET("/monitor/:monitor_id/events", types.WithEnv(monitor.MonitorEvents), IsAuthenticated)
	e.GET("/api/monitor/pause/:monitor_id", types.WithEnv(monitor.PauseMonitor), IsAuthenticated)
	e.GET("/api/monitor/resume/:monitor_id", types.WithEnv(monitor.ResumeMonitor), IsAuthenticated)
	e.POST("/api/monitors/create", types.WithEnv(monitor.CreateMonitor), IsAuthenticated)
	e.PUT("/api/monitor/:monitor_id", types.WithEnv(monitor.UpdateMonitor), IsAuthenticated)
	e.DELETE("/api/monitor/:project_id/:monitor_id", types.WithEnv(monitor.DeleteMonitor), IsAuthenticated)

	e.GET("/monitor/:monitor_id/integrations", types.WithEnv(monitor.MonitorIntegrations), IsAuthenticated)
	e.PUT("/monitor/:monitor_id/integrations", types.WithEnv(monitor.UpdateMonitorIntegrations), IsAuthenticated)

	e.GET("/projects", types.WithEnv(project.Projects), IsAuthenticated)
	e.POST("/api/projects/create", types.WithEnv(project.CreateProject), IsAuthenticated)

	e.GET("/p/:ping_slug", types.WithEnv(ping.Ping))
	e.POST("/api/contactus", types.WithEnv(dashboard.SubmitContact))
	e.POST("/process-payment", handlePayment)

	paypalClient, err = paypal.NewClient(clientID, clientSecret, paypal.APIBaseSandBox)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	_, err = paypalClient.GetAccessToken()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	l.Log.Info("Starting server")
	e.Logger.Fatal(e.Start(":1323"))

}
