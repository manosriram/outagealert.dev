package dashboard

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

type ContactForm struct {
	Name             string `form:"name"`
	Email            string `form:"email"`
	ENV              string `form:"ENV"`
	Message          string `form:"message"`
	PaymentSessionId string
}

func DashboardHome(c echo.Context) error {
	return c.Render(200, "dashboard.html", nil)
}

func Pricing(c echo.Context) error {
	return c.Render(200, "pricing.html", ContactForm{Name: "mano", ENV: os.Getenv("ENV")})
}

func Faq(c echo.Context) error {
	return c.Render(200, "faq.html", nil)
}

func Terms(c echo.Context) error {
	return c.Render(200, "termsconditions.html", nil)
}

func Contact(c echo.Context) error {
	return c.Render(200, "contact.html", nil)
}

func Refund(c echo.Context) error {
	return c.Render(200, "refund.html", nil)
}

func EmailVerified(c echo.Context) error {
	return c.Render(200, "email-verified.html", nil)
}

func SubmitContact(c echo.Context, env *types.Env) error {
	contactForm := new(ContactForm)

	if err := c.Bind(contactForm); err != nil {
		return c.Render(200, "errors", template.Response{Error: "Invalid form data"})
	}

	if contactForm.Email == "" || contactForm.Message == "" || contactForm.Name == "" {
		return c.Render(200, "errors", template.Response{Error: "Invalid form data"})
	}
	err := env.DB.Query.AddContactFormEntry(c.Request().Context(), db.AddContactFormEntryParams{Name: &contactForm.Name, Email: &contactForm.Email, Message: &contactForm.Message})
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}
	return c.Render(200, "errors", template.Response{Message: "Form submitted"})
}
