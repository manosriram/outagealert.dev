package monitor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/ping"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	PING_HOST            = "https://ping.outagealert.io"
	NANOID_ALPHABET_LIST = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"
	NANOID_LENGTH        = 22
)

type CreateMonitorForm struct {
	Name        string `form:"name" validate:"required"`
	Period      int32  `form:"period" validate:"required"`
	GracePeriod int32  `form:"grace_period" validate:"required"`
	ProjectId   string `form:"project_id" validate:"required"`
}

func getMonitorFromFetchedRow(fetchedRow db.GetAllMonitorIDsRow) db.Monitor {
	return db.Monitor{
		ID:          fetchedRow.ID,
		Period:      fetchedRow.Period,
		GracePeriod: fetchedRow.GracePeriod,
	}
}

func CreateMonitor(c echo.Context, env *types.Env) error {
	createMonitorForm := new(CreateMonitorForm)
	if err := c.Bind(createMonitorForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}
	s, err := session.Get("session", c)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}
	email := s.Values["email"].(string)

	pingSlug, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	pingUrl := fmt.Sprintf("%s/%s", PING_HOST, pingSlug)
	fmt.Println(pingUrl)

	id, err := gonanoid.New()
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	monitor, err := env.DB.Query.CreateMonitor(c.Request().Context(), db.CreateMonitorParams{
		ID:          id,
		ProjectID:   createMonitorForm.ProjectId,
		PingUrl:     pingUrl,
		Type:        "",
		UserEmail:   email,
		Name:        createMonitorForm.Name,
		Period:      createMonitorForm.Period,
		GracePeriod: createMonitorForm.GracePeriod,
	})
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: err.Error()})
	}

	go ping.StartMonitorCheck(monitor, env)

	return c.Render(200, "monitor-list-block", template.UserMonitor{Monitor: monitor})
}

func GetMonitorEvents(c echo.Context, env *types.Env) error {
	fmt.Println("hit")
	monitorId := c.Param("monitor_id")

	events, err := env.DB.Query.GetEventsByMonitorId(c.Request().Context(), monitorId)
	if err != nil {
		return err
	}

	return c.Render(200, "monitor-events", template.MonitorEvents{Events: events})
}

func StartAllMonitorChecks(env *types.Env) {
	fmt.Println("started checks")
	monitors, err := env.DB.Query.GetAllMonitorIDs(context.Background())
	if err != nil {
	}

	for _, monitor := range monitors {
		go ping.StartMonitorCheck(getMonitorFromFetchedRow(monitor), env)
	}
}
