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
	Email            string
	MonitorName      string
	MonitorId        string
	MonitorLink      string
	Env              types.Env
	MagicToken       string
	OTP              string
	NotificationType NotificationType
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
	OTP       string
}

type MonitorDownAlertMailData struct {
	Subject     string
	MonitorName string
	MonitorLink string
}

// Mail represents a email request
type Mail struct {
	from       string
	to         []string
	subject    string
	body       string
	mtype      MailType
	data       interface{}
	templateId string
}

// CreateMail takes in a mail request and constructs a sendgrid mail type.
func (e EmailNotification) CreateMail(mailReq *Mail) []byte {

	m := mail.NewV3Mail()

	from := mail.NewEmail("outagealert", mailReq.from)
	m.SetFrom(from)

	m.SetTemplateID(mailReq.templateId)

	p := mail.NewPersonalization()

	tos := make([]*mail.Email, 0)
	for _, to := range mailReq.to {
		tos = append(tos, mail.NewEmail("user", to))
	}

	p.AddTos(tos...)

	switch mailReq.data.(type) {
	case *VerifyEmailMailData:
		d := mailReq.data.(*VerifyEmailMailData)
		p.SetDynamicTemplateData("name", d.Name)
		p.SetDynamicTemplateData("host_with_scheme", os.Getenv("HOST_WITH_SCHEME"))
		p.SetDynamicTemplateData("magic_link", d.MagicLink)
		p.SetDynamicTemplateData("otp", d.OTP)
	case *MonitorDownAlertMailData:
		d := mailReq.data.(*MonitorDownAlertMailData)
		p.SetDynamicTemplateData("monitor_name", d.MonitorName)
		p.SetDynamicTemplateData("monitor_link", d.MonitorLink)
	}

	m.AddPersonalizations(p)
	return mail.GetRequestBody(m)
}

func (e EmailNotification) SendMail(mailType, templateId string, data interface{}) error {
	l.Log.Infof("Sending email to %s", e.Email)
	switch mailType {
	case string(VERIFY_EMAIL), string(FORGOT_PASSWORD_OTP):
		d := data.(VerifyEmailMailData)
		b := e.CreateMail(&Mail{
			from:       os.Getenv("SMTP_EMAIL"),
			to:         []string{e.Email},
			subject:    d.Subject,
			data:       &d,
			templateId: templateId,
		})
		return e.DeliverMail(b)
	case string(MONITOR_DOWN), string(MONITOR_UP):
		d := data.(MonitorDownAlertMailData)
		b := e.CreateMail(&Mail{
			from:       os.Getenv("SMTP_EMAIL"),
			to:         []string{e.Email},
			subject:    d.Subject,
			data:       &d,
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

	if response.StatusCode/100 != 2 {
		l.Log.Infof("Not able to send mail %d", response.Body)
	}
	return nil
}

func (e EmailNotification) Notify() error {
	var templateId SendGridTemplateId = NotificationTypeVsTemplateId[e.NotificationType]
	go e.SendMail(string(e.NotificationType), string(templateId), MonitorDownAlertMailData{
		MonitorName: e.MonitorName,
		MonitorLink: e.MonitorLink,
	})

	return nil
}

func (e EmailNotification) SendAlert() error {
	emailIntegration, _ := e.Env.DB.Query.GetMonitorIntegration(context.Background(), db.GetMonitorIntegrationParams{
		MonitorID: e.MonitorId,
		AlertType: "email",
	})
	if !emailIntegration.EmailAlertSent {
		l.Log.Infof("Sending email alert to %s", e.MonitorId)
		err := e.Notify()
		if err != nil {
			l.Log.Errorf("Error notifying via email alert, monitor_id %s, err %s", e.MonitorId, err.Error())
			return err
		}

		// only mark notification as sent if integration.NotificationVsShouldMarkEmailSent[e.EmailNotificationType] is true
		if NotificationVsShouldMarkNotificationSent[e.NotificationType] {
			err = e.Env.DB.Query.UpdateEmailAlertSentFlag(context.Background(), db.UpdateEmailAlertSentFlagParams{
				MonitorID:      e.MonitorId,
				EmailAlertSent: true,
			})
			if err != nil {
				l.Log.Errorf("Error updating email alert sent flag %s", err.Error())
				return err
			}
		}
	}

	return nil
}
