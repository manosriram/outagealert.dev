package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/models"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

type Env struct {
	Users     models.UserModel
	Validator *validator.Validate
}

func NewEnv(conn *db.Queries) *Env {
	return &Env{
		Users:     models.UserModel{Db: conn},
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func InjectEnv(env *Env) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("env", env)
			return next(c)
		}
	}
}

func WithEnv(h func(echo.Context, *Env) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(c, c.Get("env").(*Env))
	}
}
