package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	"github.com/rs/zerolog/log"
)

type WebhookNotification struct {
	Url         string
	MonitorName string
	MonitorId   string
	Env         types.Env
}

func (w WebhookNotification) Notify() error {
	hookUrl, err := url.Parse(w.Url)
	if err != nil {
		log.Error().Msgf("Error getting url %s", hookUrl)
		return err
	}

	response, err := http.Get(hookUrl.String())
	if err != nil {
		log.Error().Msgf("Error getting url %s", hookUrl)
		return err
	}
	fmt.Println(response)
	return nil
}

func (w WebhookNotification) SendAlert(monitorId, monitorName string) error {
	w.MonitorName = monitorName
	w.MonitorId = monitorId
	integs, err := w.Env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: monitorId,
		AlertType: "webhook",
	})
	if err != nil {
		log.Error().Msgf("Error sending email alert, monitor_id %s, err %s", monitorId, err.Error())
		return err
	}
	if !integs.WebhookAlertSent {
		err := w.Notify()
		if err != nil {
			return err
		}
		w.Env.DB.Query.UpdateWebhookAlertSentFlag(context.Background(), db.UpdateWebhookAlertSentFlagParams{
			MonitorID:        monitorId,
			WebhookAlertSent: true,
		})
	}

	return nil
}
