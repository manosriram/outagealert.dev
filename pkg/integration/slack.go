package integration

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/types"
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

	s, _ := session.Get("session", c)
	email := s.Values["email"].(string)
	host := os.Getenv("HOST_WITH_SCHEME")

	fmt.Println(email, host)

	return nil
}
