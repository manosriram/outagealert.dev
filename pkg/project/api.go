package project

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

type CreateProjectForm struct {
	Name       string `form:"name" validate:"required"`
	Visibility string `form:"visibility"`
}

func CreateProject(c echo.Context, env *types.Env) error {
	createProjectForm := new(CreateProjectForm)
	if err := c.Bind(createProjectForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}
	s, _ := session.Get("session", c)
	email := s.Values["email"].(string)

	err := env.DB.Query.CreateProject(c.Request().Context(), db.CreateProjectParams{Name: createProjectForm.Name, Visibility: createProjectForm.Visibility, UserEmail: email})
	fmt.Println("created ", err)

	return nil
}
