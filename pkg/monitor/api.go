package monitor

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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

type UpdateMonitorForm struct {
	Name        string `form:"name" validate:"required"`
	Period      int32  `form:"period" validate:"required"`
	GracePeriod int32  `form:"grace_period" validate:"required"`
	MonitorId   string `form:"monitor_id" validate:"required"`
}

type DeleteMonitorForm struct {
	MonitorId string `form:"monitor_id" validate:"required"`
}

func getMonitorFromFetchedRow(fetchedRow db.GetAllMonitorIDsRow) db.Monitor {
	return db.Monitor{
		ID:          fetchedRow.ID,
		Period:      fetchedRow.Period,
		GracePeriod: fetchedRow.GracePeriod,
	}
}

func DeleteMonitor(c echo.Context, env *types.Env) error {
	s, err := session.Get("session", c)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}
	email := s.Values["email"].(string)
	monitorId := c.Param("monitor_id")

	err = env.DB.Query.DeleteMonitor(c.Request().Context(), db.DeleteMonitorParams{
		ID:        monitorId,
		UserEmail: email,
	})
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	c.Request().Header.Set("HX-Redirect", "/projects")
	return c.NoContent(200)
	// return c.Render(200, "monitors.html", template.Response{Message: "Monitor deleted successfully"})
}

func UpdateMonitor(c echo.Context, env *types.Env) error {
	updateMonitorForm := new(UpdateMonitorForm)
	if err := c.Bind(updateMonitorForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}
	s, err := session.Get("session", c)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}
	email := s.Values["email"].(string)
	monitorId := c.Param("monitor_id")

	err = env.DB.Query.UpdateMonitor(c.Request().Context(), db.UpdateMonitorParams{
		Name:        updateMonitorForm.Name,
		Period:      updateMonitorForm.Period,
		GracePeriod: updateMonitorForm.GracePeriod,
		ID:          monitorId,
		UserEmail:   email,
	})
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	return c.Render(200, "errors", template.Response{Message: "Monitor updated successfully", Error: "This is err"})
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

func GetMonitorEventsTable(c echo.Context, env *types.Env) error {
	monitorId := c.Param("monitor_id")
	page := c.QueryParam("page")
	pageNumber, err := strconv.Atoi(page)
	if pageNumber <= 0 {
		pageNumber = 1
	}
	offset := (pageNumber - 1) * 25

	events, err := env.DB.Query.GetEventsByMonitorIdPaginated(c.Request().Context(), db.GetEventsByMonitorIdPaginatedParams{
		MonitorID: monitorId,
		Offset:    int32(offset),
	})
	if err != nil {
		fmt.Println("e = ", err)
		return err
	}

	hasNextPage := true
	if len(events) == 0 {
		hasNextPage = false
	}
	return c.Render(200, "monitor-events-table", template.MonitorEvents{MonitorID: monitorId, Events: events, CurrentPage: pageNumber, NextPage: pageNumber + 1, HasNextPage: hasNextPage})
}

func GetMonitorEvents(c echo.Context, env *types.Env) error {
	fmt.Println("hit")
	monitorId := c.Param("monitor_id")
	page := c.QueryParam("page")
	pageInt, err := strconv.Atoi(page)
	offset := (pageInt - 1) * 25
	fmt.Println(err)

	events, err := env.DB.Query.GetEventsByMonitorIdPaginated(c.Request().Context(), db.GetEventsByMonitorIdPaginatedParams{
		MonitorID: monitorId,
		Offset:    int32(offset),
	})
	if err != nil {
		fmt.Println("e = ", err)
		return err
	}

	hasNextPage := true
	if len(events) == 0 {
		hasNextPage = false
	}
	return c.Render(200, "monitor-events", template.MonitorEvents{MonitorID: monitorId, Events: events, CurrentPage: pageInt, NextPage: pageInt + 1, HasNextPage: hasNextPage})
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
