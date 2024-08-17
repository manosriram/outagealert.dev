package monitor

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

const (
	PING_HOST = "https://ping.outagealert.io"
)

type CreateMonitorForm struct {
	Name        string `form:"name" validate:"required"`
	Period      int32  `form:"period" validate:"required"`
	GracePeriod int32  `form:"grace_period" validate:"required"`
	ProjectId   string `form:"project_id" validate:"required"`
}

func CreateMonitor(c echo.Context, env *types.Env) error {
	createMonitorForm := new(CreateMonitorForm)
	if err := c.Bind(createMonitorForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}
	s, _ := session.Get("session", c)
	email := s.Values["email"].(string)
	project_uuid, _ := uuid.Parse(createMonitorForm.ProjectId)

	pingSlug := uuid.New()
	pingUrl := fmt.Sprintf("%s/%s", PING_HOST, pingSlug)

	err := env.DB.Query.CreateMonitor(c.Request().Context(), db.CreateMonitorParams{ProjectID: pgtype.UUID{Bytes: project_uuid, Valid: true}, PingUrl: pingUrl, Type: pgtype.Text{}, UserEmail: email, Name: createMonitorForm.Name, Period: createMonitorForm.Period, GracePeriod: createMonitorForm.GracePeriod})
	fmt.Println("created ", err)

	return nil
}
