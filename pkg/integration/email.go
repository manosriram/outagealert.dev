package integration

import (
	"context"
	"os"

	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailNotification struct {
	Email       string
	MonitorName string
	MonitorId   string
	Env         types.Env
	MagicToken  string
	OTP         string
}

type MailType int

// List of Mail Types we are going to send.
const (
	MailConfirmation MailType = iota + 1
	PassReset
)

// VerifyEmailMailData represents the data to be sent to the template of the mail.
type VerifyEmailMailData struct {
	Name      string
	Subject   string
	MagicLink string
	Host      string
	OTP       string
}

// Mail represents a email request
type Mail struct {
	from       string
	to         []string
	subject    string
	body       string
	mtype      MailType
	data       *VerifyEmailMailData
	templateId string
}

// CreateMail takes in a mail request and constructs a sendgrid mail type.
func (e EmailNotification) CreateMail(mailReq *Mail) []byte {

	m := mail.NewV3Mail()

	from := mail.NewEmail("Mano Sriram", mailReq.from)
	m.SetFrom(from)

	m.SetTemplateID(mailReq.templateId)

	p := mail.NewPersonalization()

	tos := make([]*mail.Email, 0)
	for _, to := range mailReq.to {
		tos = append(tos, mail.NewEmail("user", to))
	}

	p.AddTos(tos...)

	p.SetDynamicTemplateData("name", mailReq.data.Name)
	p.SetDynamicTemplateData("host", mailReq.data.Host)
	p.SetDynamicTemplateData("magic_link", mailReq.data.MagicLink)
	p.SetDynamicTemplateData("otp", mailReq.data.OTP)

	m.AddPersonalizations(p)
	return mail.GetRequestBody(m)
}

func (e EmailNotification) SendMail(mailType, templateId string, data VerifyEmailMailData) error {
	switch mailType {
	case "verify_email":
		b := e.CreateMail(&Mail{
			from:       os.Getenv("SMTP_EMAIL"),
			to:         []string{e.Email},
			subject:    data.Subject,
			data:       &data,
			templateId: templateId,
		})
		return e.DeliverMail(b)
	case "forgot_password_otp":
		b := e.CreateMail(&Mail{
			from:       os.Getenv("SMTP_EMAIL"),
			to:         []string{e.Email},
			subject:    data.Subject,
			data:       &data,
			templateId: templateId,
		})
		return e.DeliverMail(b)
	}

	return nil
}

func (e EmailNotification) DeliverMail(body []byte) error {
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = body
	response, err := sendgrid.API(request)
	if err != nil {
		l.Log.Errorf("Unable to send mail %s", err)
		return err
	}
	l.Log.Infof("Mail sent successfully to %d", response.StatusCode)
	return nil
}

func (e EmailNotification) Notify() error {
	// subject := fmt.Sprintf("Monitor DOWN alert")
	// body := fmt.Sprintf("%s is DOWN!", e.MonitorName)
	// err := e.SendMail("", subject, body)
	// if err != nil {
	// log.Error().Msgf("Error notifying via email %s", err.Error())
	// return err
	// }
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
		l.Log.Errorf("Error sending email alert, monitor_id %s, err %s", monitorId, err.Error())
		return err
	}
	if !integs.EmailAlertSent {
		err := e.Notify()
		if err != nil {
			l.Log.Errorf("Error notifying via email alert, monitor_id %s, err %s", monitorId, err.Error())
			return err
		}
		e.Env.DB.Query.UpdateEmailAlertSentFlag(context.Background(), db.UpdateEmailAlertSentFlagParams{
			MonitorID:      monitorId,
			EmailAlertSent: true,
		})
	}

	return nil
}
