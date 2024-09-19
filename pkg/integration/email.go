package integration

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	"github.com/rs/zerolog/log"
)

type EmailNotification struct {
	Email string
	Env   types.Env
}

func (e EmailNotification) SendAlert(monitorId, monitorName string) error {
	integs, err := e.Env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: monitorId,
		AlertType: "email",
	})
	if err != nil {
		log.Error().Msgf("Error sending email alert, monitor_id %s, err %s", monitorId, err.Error())
		return err
	}
	if !integs.EmailAlertSent {
		auth := smtp.PlainAuth(
			"",
			os.Getenv("SMTP_EMAIL"),
			"rgrnnvzqezcraobh",
			"smtp.gmail.com",
		)
		msg := fmt.Sprintf("This is subject\n%s is DOWN!", monitorName)
		smtp.SendMail("smtp.gmail.com:587", auth, "mano.sriram0@gmail.com", []string{e.Email}, []byte(msg))

		e.Env.DB.Query.UpdateEmailAlertSentFlag(context.Background(), db.UpdateEmailAlertSentFlagParams{
			MonitorID:      monitorId,
			EmailAlertSent: true,
		})
	}

	return nil
}
