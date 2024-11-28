package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

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

func (s SlackNotification) Notify() error {
	fmt.Println("notifying via slack")
	// slackUser, err := s.Env.DB.Query.GetSlackUserByEmail(context.Background(), s.UserEmail)
	// if err != nil {
	// fmt.Println(err)
	// }

	var te, alertText string
	switch s.NotificationType {
	case MONITOR_DOWN:
		alertText = fmt.Sprintf(`Monitor **DOWN** alert`, s.MonitorName)
		te = "Monitor DOWN"
	case MONITOR_UP:
		alertText = fmt.Sprintf(`Monitor **UP** alert`, s.MonitorName)
		te = "Monitor UP"
	}
	url := "https://hooks.slack.com/services/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"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
  "channel": "C01H43ZPFRC",
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
          "text": "%s"
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
}`, alertText, te, s.MonitorName, string(s.NotificationType), time.Now().Format("2006-01-02150405"), s.MonitorId))

	fmt.Println(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Authorization", "Bearer it1nhwmqdiyiukk7jxbj4dczxr")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
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
		s.Env.DB.Query.UpdateSlackAlertSentFlag(context.Background(), db.UpdateSlackAlertSentFlagParams{
			MonitorID:      s.MonitorId,
			SlackAlertSent: true,
		})
	}

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
		"https://e0c2-2405-201-e07a-e037-2c3e-4f2f-7b3-dae.ngrok-free.app/integration",
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
