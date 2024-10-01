package template

import (
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/sqlc/db"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Templates struct {
	templates *template.Template
}

func Title(s db.AlertType) string {
	return cases.Title(language.English, cases.Compact).String(string(s))
}

func FormatTimeWithoutAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration.Seconds() < 60:
		return fmt.Sprintf("%d seconds", int(duration.Seconds()))
	case duration.Minutes() < 60:
		return fmt.Sprintf("%d minutes", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%d hours", int(duration.Hours()))
	case duration.Hours() < 48:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days", int(duration.Hours()/24))
	}
}

func FormatTimeWithAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration.Seconds() < 60:
		return "just now"
	case duration.Minutes() < 60:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration.Hours() < 48:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	funcs := template.FuncMap{
		"createdAtDistanceWithAgo":    FormatTimeWithAgo,
		"createdAtDistanceWithoutAgo": FormatTimeWithoutAgo,
		"title":                       Title,
	}
	return &Templates{
		templates: template.Must(template.New("").Funcs(funcs).ParseGlob("views/*.html")),
	}
}
