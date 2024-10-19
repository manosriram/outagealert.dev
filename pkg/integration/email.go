package integration

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailNotification struct {
	Email       string
	MonitorName string
	MonitorId   string
	Env         types.Env
}

func (e EmailNotification) Notify() error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
		"mail.privateemail.com",
	)
	subject := fmt.Sprintf("Monitor DOWN alert")
	body := fmt.Sprintf("%s is DOWN!", e.MonitorName)
	err := smtp.SendMail("mail.privateemail.com:587", auth, "hello@outagealert.dev", []string{e.Email}, []byte(subject))
	if err != nil {
		fmt.Println("error sending mail ", err)
	}

	from := mail.NewEmail("Mano Sriram", os.Getenv("SMTP_EMAIL"))
	to := mail.NewEmail("", e.Email)
	plainTextContent := body
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, body)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}

func (e EmailNotification) SendAlert(monitorId, monitorName string) error {
	e.MonitorName = monitorName
	e.MonitorId = monitorId
	integs, err := e.Env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: monitorId,
		AlertType: "email",
	})
	if err != nil {
		log.Error().Msgf("Error sending email alert, monitor_id %s, err %s", monitorId, err.Error())
		return err
	}
	if !integs.EmailAlertSent {
		err := e.Notify()
		if err != nil {
			log.Error().Msgf("Error notifying via email alert, monitor_id %s, err %s", monitorId, err.Error())
			return err
		}
		e.Env.DB.Query.UpdateEmailAlertSentFlag(context.Background(), db.UpdateEmailAlertSentFlagParams{
			MonitorID:      monitorId,
			EmailAlertSent: true,
		})
	}

	return nil
}
