package integration

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
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
	// s, _ := session.Get("session", c)
	// email := s.Values["email"].(string)

	fmt.Println("handle slack auth")
	code := c.QueryParam("code")
	email := c.QueryParam("state")
	resp, err := slack.GetOAuthV2Response(
		http.DefaultClient,
		os.Getenv("SLACK_CLIENT_ID"),
		os.Getenv("SLACK_CLIENT_SECRET"),
		code,
		"https://94fa-2405-201-e07a-e037-71fa-143e-b6b5-5520.ngrok-free.app/integration",
	)
	if err != nil {
		l.Log.Errorf("Error getting oauth v2 response %s", err.Error())
		return c.JSON(500, nil)
	}

	slackUser, err := env.DB.Query.GetSlackUserByEmail(c.Request().Context(), email)
	if slackUser.ChannelName != nil {
		err = env.DB.Query.UpdateSlackUserByEmail(c.Request().Context(), db.UpdateSlackUserByEmailParams{
			UserEmail:        email,
			ChannelUrl:       &resp.IncomingWebhook.URL,
			ChannelName:      &resp.IncomingWebhook.Channel,
			ChannelID:        &resp.IncomingWebhook.ChannelID,
			ConfigurationUrl: &resp.IncomingWebhook.ConfigurationURL,
		})
		if err != nil {
			l.Log.Errorf("Error updating slack user %s", err.Error())
			return c.JSON(500, nil)
		}
	}

	err = env.DB.Query.CreateNewSlackUser(c.Request().Context(), db.CreateNewSlackUserParams{
		UserEmail:        email,
		ChannelUrl:       &resp.IncomingWebhook.URL,
		ChannelName:      &resp.IncomingWebhook.Channel,
		ChannelID:        &resp.IncomingWebhook.ChannelID,
		ConfigurationUrl: &resp.IncomingWebhook.ConfigurationURL,
	})
	if err != nil {
		l.Log.Errorf("Error creating new slack user %s", err.Error())
		return c.JSON(500, nil)
	}

	return c.JSON(200, nil)
}
