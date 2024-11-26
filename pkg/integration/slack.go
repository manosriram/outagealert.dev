package integration

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/slack-go/slack"
)

type SlackNotification struct {
	MonitorName      string
	MonitorId        string
	MonitorLink      string
	Env              types.Env
	NotificationType NotificationType
}

func (s SlackNotification) Notify() error {
	return nil
}

func (s SlackNotification) SendAlert(monitorId, monitorName string) error {
	return nil
}

type SlackAuthPayload struct {
}

func HandleSlackAuth(c echo.Context, env *types.Env) error {
	slackAuthPayload := new(SlackAuthPayload)
	if err := c.Bind(slackAuthPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	// s, _ := session.Get("session", c)
	// email := s.Values["email"].(string)
	// host := os.Getenv("HOST_WITH_SCHEME")

	code := c.QueryParam("code")
	resp, err := slack.GetOAuthV2Response(
		http.DefaultClient,
		os.Getenv("SLACK_CLIENT_ID"),
		os.Getenv("SLACK_CLIENT_SECRET"),
		code,
		"https://0183-2405-201-e07a-e037-2c3e-4f2f-7b3-dae.ngrok-free.app",
	)

	fmt.Println(resp, err)
	return nil
}
