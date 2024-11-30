package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

type WebhookNotification struct {
	Url              string
	MonitorName      string
	MonitorId        string
	Env              types.Env
	NotificationType NotificationType
}

func (w WebhookNotification) Notify() error {
	hookUrl, err := url.Parse(w.Url)
	if err != nil {
		l.Log.Errorf("Error getting url %s", hookUrl)
		return err
	}

	response, err := http.Get(hookUrl.String())
	if err != nil {
		// l.Log.Errorf("Error getting url %s", hookUrl)
		return err
	}
	fmt.Println(response)
	return nil
}

func (w WebhookNotification) SendAlert() error {
	integs, err := w.Env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: w.MonitorId,
		AlertType: "webhook",
	})
	if err != nil {
		l.Log.Errorf("Error sending email alert, monitor_id %s, err %s", w.MonitorId, err.Error())
		return err
	}
	if !integs.WebhookAlertSent {
		err := w.Notify()
		if err != nil {
			return err
		}

		if NotificationVsShouldMarkNotificationSent[w.NotificationType] {
			w.Env.DB.Query.UpdateWebhookAlertSentFlag(context.Background(), db.UpdateWebhookAlertSentFlagParams{
				MonitorID:        w.MonitorId,
				WebhookAlertSent: true,
			})
		}
	}

	return nil
}
