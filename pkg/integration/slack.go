package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/session"
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
	UserEmail        string
	Env              types.Env
	NotificationType NotificationType
}

type Message struct {
	Text   string       `json:"text"`
	Blocks []SlackBlock `json:"blocks"`
}

type SlackBlock struct {
	Type string  `json:"type"`
	Text Section `json:"text"`
}

type Section struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var monitorAlertTemplate string = `{
  "channel": "%s",
  "text": "%s",
  "blocks": [
    {
      "type": "header",
      "text": {
      "type": "plain_text",
        "text": "%s",
        "emoji": true
      }
    },
    {
      "type": "section",
      "fields": [
        {
          "type": "mrkdwn",
          "text": "*Monitor Name:*\n%s"
        },
        {
          "type": "mrkdwn",
          "text": "*Status:*\n%s"
        }
      ]
    },
    {
      "type": "section",
      "fields": [
        {
          "type": "mrkdwn",
				"text": "*Timestamp:*\n%s"
        }
      ]
    },
    {
      "type": "actions",
      "elements": [
        {
          "type": "button",
          "text": {
            "type": "plain_text",
            "emoji": true,
            "text": "Visit monitor"
          },
          "style": "primary",
					"url": "%s"
        }
      ]
    }
  ]
}`

func (s SlackNotification) Notify() error {
	slackUser, err := s.Env.DB.Query.GetSlackUserByEmail(context.Background(), s.UserEmail)
	if err != nil {
		l.Log.Errorf("Error getting slack user by email %s", err.Error())
		return err
	}

	var alertText, alertDescription string
	switch s.NotificationType {
	case MONITOR_DOWN:
		alertDescription = fmt.Sprintf(`Monitor **DOWN** alert`, s.MonitorName)
		alertText = "Monitor DOWN"
	case MONITOR_UP:
		alertDescription = fmt.Sprintf(`Monitor **UP** alert`, s.MonitorName)
		alertText = "Monitor UP"
	}
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(monitorAlertTemplate, slackUser.ChannelID, alertDescription, alertText, s.MonitorName, string(s.NotificationType), time.Now().Format("2006-01-02-15:04:05"), s.MonitorLink))

	client := &http.Client{}
	req, err := http.NewRequest(method, *slackUser.ChannelUrl, payload)

	if err != nil {
		l.Log.Errorf("Error making slack HTTP request %s", err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		l.Log.Errorf("Error making slack HTTP request %s", err.Error())
		return err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		l.Log.Errorf("Error reading response body %s")
		return err
	}
	return nil
}

func (s SlackNotification) SendAlert() error {
	integs, err := s.Env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: s.MonitorId,
		AlertType: "slack",
	})
	if err != nil {
		l.Log.Errorf("Error sending slack alert, monitor_id %s, err %s", s.MonitorId, err.Error())
		return err
	}
	if !integs.SlackAlertSent {
		err := s.Notify()
		if err != nil {
			return err
		}

		if NotificationVsShouldMarkNotificationSent[s.NotificationType] {
			s.Env.DB.Query.UpdateSlackAlertSentFlag(context.Background(), db.UpdateSlackAlertSentFlagParams{
				MonitorID:      s.MonitorId,
				SlackAlertSent: true,
			})
		}
	}

	return nil
}

func DisconnectIntegration(c echo.Context, env *types.Env) error {
	s, _ := session.Get("session", c)
	email := s.Values["email"].(string)

	monitorId := c.Param("monitor_id")
	provider := c.QueryParam("provider")

	return DisconnectProvider{
		C:         c,
		Env:       env,
		Email:     email,
		MonitorId: monitorId,
		Provider:  provider,
	}.Disconnect()
}

func HandleSlackAuth(c echo.Context, env *types.Env) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	state, err := Base64Decode(state)
	if err != nil {
		l.Log.Errorf("Invalid base64 encoded state %s", state)
		return c.JSON(500, nil)
	}

	decodedState := strings.Split(state, ";")
	if len(decodedState) != 2 {
		l.Log.Errorf("Invalid base64 decoded state %s", decodedState)
		return c.JSON(500, nil)
	}

	projectId := decodedState[0]
	monitorId := decodedState[1]

	monitor, err := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	if monitor.ID == "" {
		l.Log.Errorf("Error getting monitor %s", err.Error())
		return c.JSON(500, nil)
	}

	resp, err := slack.GetOAuthV2Response(
		http.DefaultClient,
		os.Getenv("SLACK_CLIENT_ID"),
		os.Getenv("SLACK_CLIENT_SECRET"),
		code,
		os.Getenv("SLACK_REDIRECT_URL"),
	)
	if err != nil {
		l.Log.Errorf("Error getting oauth v2 response %s", err.Error())
		return c.JSON(500, nil)
	}

	slackUser, err := env.DB.Query.GetSlackUserByEmail(c.Request().Context(), monitor.UserEmail)
	if slackUser.ChannelName != nil {
		err = env.DB.Query.UpdateSlackUserByEmail(c.Request().Context(), db.UpdateSlackUserByEmailParams{
			UserEmail:        monitor.UserEmail,
			ChannelUrl:       &resp.IncomingWebhook.URL,
			ChannelName:      &resp.IncomingWebhook.Channel,
			ChannelID:        &resp.IncomingWebhook.ChannelID,
			ConfigurationUrl: &resp.IncomingWebhook.ConfigurationURL,
		})
		if err != nil {
			l.Log.Errorf("Error updating slack user %s", err.Error())
			return c.JSON(500, nil)
		}
	} else {
		err = env.DB.Query.CreateNewSlackUser(c.Request().Context(), db.CreateNewSlackUserParams{
			UserEmail:        monitor.UserEmail,
			ChannelUrl:       &resp.IncomingWebhook.URL,
			ChannelName:      &resp.IncomingWebhook.Channel,
			ChannelID:        &resp.IncomingWebhook.ChannelID,
			ConfigurationUrl: &resp.IncomingWebhook.ConfigurationURL,
		})
		if err != nil {
			l.Log.Errorf("Error creating new slack user %s", err.Error())
			return c.JSON(500, nil)
		}
	}

	err = env.DB.Query.UpdateSlackAlertIntegration(c.Request().Context(), db.UpdateSlackAlertIntegrationParams{
		MonitorID: monitorId,
		IsActive:  true,
	})
	if err != nil {
		l.Log.Errorf("Error updating webhook alert integration %s", err.Error())
		return c.JSON(500, nil)
	}

	return c.Redirect(301, makeSlackRedirectUrl(projectId, monitorId))
}
