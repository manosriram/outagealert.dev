package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

type SigninForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type SignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func SignInApi(c echo.Context, env *types.Env) error {
	signinForm := new(SigninForm)
	if err := c.Bind(signinForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Println(err)
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["email"] = signinForm.Email
	c.Response().Header().Set("HX-Redirect", "/dashboard")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		fmt.Println("err ", err.Error())
		return err
	}
	return c.Render(200, "signup-success.html", template.RegisterSuccessResponse{Email: signinForm.Email})
}

func SignUpApi(c echo.Context, env *types.Env) error {
	signupForm := new(SignupForm)
	if err := c.Bind(signupForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}
	fmt.Println(signupForm.Email, signupForm.Password, signupForm.Name)
	err := env.Users.Create(&db.User{
		Email:    signupForm.Email,
		Password: pgtype.Text{String: signupForm.Password, Valid: true},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return c.Render(200, "signup.html", template.Response{Message: "notok", Error: "User already exists"})
			default:
				return c.Render(200, "signup.html", template.Response{Message: "notok", Error: "Internal server error"})
			}
		}
	}
	return c.Render(200, "signup-success.html", template.RegisterSuccessResponse{Email: signupForm.Email})
}
