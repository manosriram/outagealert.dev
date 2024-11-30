package integration

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

type DisconnectProvider struct {
	C         echo.Context
	Env       *types.Env
	Email     string
	MonitorId string
	Provider  string
}

func (d DisconnectProvider) Disconnect() error {
	d.C.Response().Header().Set("HX-Refresh", "true")
	switch d.Provider {
	case "slack":
		err := d.DisconnectSlack()
		if err != nil {
			l.Log.Errorf("Error disconnecting slack %s", err.Error())
			return d.C.Render(400, "errors", template.Response{Error: "Error disconnecting slack provider"})
		}
	default:
		return d.C.Render(400, "errors", template.Response{Error: "Unknown integration"})
	}
	return d.C.Render(200, "errors", template.Response{Message: fmt.Sprintf("Disconnected %s integration", d.Provider)})
}

func (d DisconnectProvider) DisconnectSlack() error {
	err := d.Env.DB.Query.DeleteSlackUserByEmail(d.C.Request().Context(), d.Email)
	if err != nil {
		d.C.Response().Header().Set("HX-Retarget", "#error-container")
		l.Log.Errorf("Error deleting slack user by email %s", err.Error())
		return err
	}
	err = d.Env.DB.Query.UpdateSlackAlertIntegration(d.C.Request().Context(), db.UpdateSlackAlertIntegrationParams{
		MonitorID: d.MonitorId,
		IsActive:  false,
	})
	if err != nil {
		l.Log.Errorf("Error updating webhook alert integration %s", err.Error())
		return err
	}
	return d.C.Render(200, "errors", template.Response{Message: fmt.Sprintf("Disconnected %s integration", d.Provider)})
}
