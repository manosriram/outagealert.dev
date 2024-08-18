package project

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	NANOID_ALPHABET_LIST = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"
	NANOID_LENGTH        = 22
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

	id, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	project, err := env.DB.Query.CreateProject(c.Request().Context(), db.CreateProjectParams{ID: id, Name: createProjectForm.Name, Visibility: createProjectForm.Visibility, UserEmail: email})
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	return c.Render(200, "project-list-block", template.UserProject{Project: project})
}
