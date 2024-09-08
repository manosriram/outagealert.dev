package monitor

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/event"
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

func ResumeMonitor(c echo.Context, env *types.Env) error {
	monitorId := c.Param("monitor_id")

	oldMonitor, err := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	timeDiffBetweenPauseAndResume := time.Now().UTC().Sub(oldMonitor.LastPausedAt.Time)
	var newTotalPauseTime int32
	if oldMonitor.TotalPauseTime != nil {
		newTotalPauseTime = *oldMonitor.TotalPauseTime + int32(timeDiffBetweenPauseAndResume.Seconds())
	} else {
		newTotalPauseTime = int32(timeDiffBetweenPauseAndResume.Seconds())
	}
	fmt.Println("total pause time = ", newTotalPauseTime)

	updatedMonitor, err := env.DB.Query.ResumeMonitor(c.Request().Context(), db.ResumeMonitorParams{
		ID:             monitorId,
		LastResumedAt:  pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		TotalPauseTime: &newTotalPauseTime,
	})
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	oldMonitor, _ = env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	fmt.Println("last ping after resume ", oldMonitor.LastPing.Time)

	err = event.CreateEvent(context.Background(), monitorId, "paused", updatedMonitor.Status, env)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	return c.Render(200, "monitor-options", template.UserMonitor{
		Monitor: updatedMonitor,
	})
}

func PauseMonitor(c echo.Context, env *types.Env) error {
	monitorId := c.Param("monitor_id")
	oldMonitor, err := env.DB.Query.GetMonitorById(c.Request().Context(), monitorId)
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	updatedMonitor, err := env.DB.Query.PauseMonitor(c.Request().Context(), db.PauseMonitorParams{
		ID:                monitorId,
		Status:            "paused",
		StatusBeforePause: &oldMonitor.Status,
		LastPausedAt:      pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	err = event.CreateEvent(context.Background(), monitorId, oldMonitor.Status, "paused", env)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	return c.Render(200, "monitor-options", template.UserMonitor{
		Monitor: updatedMonitor,
	})
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

	c.Response().Header().Set("HX-Redirect", "/projects")
	return c.NoContent(200)
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

	// pingUrl := fmt.Sprintf("%s/%s", PING_HOST, pingSlug)
	// fmt.Println(pingUrl)

	id, err := gonanoid.New()
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	monitor, err := env.DB.Query.CreateMonitor(c.Request().Context(), db.CreateMonitorParams{
		ID:          id,
		ProjectID:   createMonitorForm.ProjectId,
		PingUrl:     pingSlug,
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

	err = event.CreateEvent(c.Request().Context(), id, "created", "up", env)
	if err != nil {
		fmt.Println(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
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

	activity, err := env.DB.Query.GetMonitorActivityPaginated(c.Request().Context(), db.GetMonitorActivityPaginatedParams{
		MonitorID: monitorId,
		Offset:    int32(offset),
	})
	if err != nil {
		fmt.Println("e = ", err)
		return err
	}

	hasNextPage := true
	if len(activity) == 0 {
		hasNextPage = false
	}
	return c.Render(200, "monitor-events-table", template.MonitorEvents{MonitorID: monitorId, Activity: activity, CurrentPage: pageNumber, NextPage: pageNumber + 1, HasNextPage: hasNextPage})
}

func GetMonitorActivity(c echo.Context, env *types.Env) error {
	monitorId := c.Param("monitor_id")
	page := c.QueryParam("page")
	pageInt, _ := strconv.Atoi(page)
	offset := (pageInt - 1) * 25

	activities, err := env.DB.Query.GetMonitorActivityPaginated(c.Request().Context(), db.GetMonitorActivityPaginatedParams{
		MonitorID: monitorId,
		Offset:    int32(offset),
	})
	// for _, activity := range activities {
	// // activity.CreatedAt = dis
	// }

	// events, err := env.DB.Query.GetEventsByMonitorIdPaginated(c.Request().Context(), db.GetEventsByMonitorIdPaginatedParams{
	// MonitorID: monitorId,
	// Offset:    int32(offset),
	// })
	if err != nil {
		fmt.Println("e = ", err)
		return err
	}

	hasNextPage := true
	if len(activities) == 0 {
		hasNextPage = false
	}
	return c.Render(200, "monitor-events", template.MonitorEvents{MonitorID: monitorId, Activity: activities, CurrentPage: pageInt, NextPage: pageInt + 1, HasNextPage: hasNextPage})
}

func StartAllMonitorChecks(env *types.Env) {
	fmt.Println("started checks")
	monitors, err := env.DB.Query.GetAllMonitorIDs(context.Background())
	if err != nil {
		fmt.Println("err ", err)
	}

	for _, monitor := range monitors {
		go ping.StartMonitorCheck(getMonitorFromFetchedRow(monitor), env)
	}
}
