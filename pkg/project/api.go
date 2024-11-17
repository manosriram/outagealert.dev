package project

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/l"
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
	host := os.Getenv("HOST_WITH_SCHEME")

	id, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
	if err != nil {
		l.Log.Error(err.Error())
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/projects", host))
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	_, err = env.DB.Query.CreateProject(c.Request().Context(), db.CreateProjectParams{ID: id, Name: createProjectForm.Name, Visibility: createProjectForm.Visibility, UserEmail: email})
	if err != nil {
		l.Log.Error(err.Error())
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/projects", host))
		return c.Render(200, "projects.html", template.Response{Error: "Internal server error"})
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/projects", host))
	return c.Render(200, "projects.html", template.UserProjects{})
}

func DeleteProject(c echo.Context, env *types.Env) error {
	projectId := c.Param("project_id")
	host := os.Getenv("HOST_WITH_SCHEME")

	err := env.DB.Query.DeleteProject(c.Request().Context(), projectId)
	if err != nil {
		l.Log.Error(err.Error())
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/projects", host))
		return c.Render(200, "projects.html", template.Response{Error: "Internal server error"})
	}

	err = env.DB.Query.DeleteProjectMonitors(c.Request().Context(), projectId)
	if err != nil {
		l.Log.Error(err.Error())
		c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/projects", host))
		return c.Render(200, "projects.html", template.Response{Error: "Internal server error"})
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/projects", host))
	return c.Render(200, "projects.html", template.UserProjects{})
}
