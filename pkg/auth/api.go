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
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

func getErrorStringFromPgxError(err error) string {
	if err != nil {
		fmt.Println("err = ", err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println("code = ", pgErr.Code)
			switch pgErr.Code {
			case "23505":
				return "User already exists"
			default:
				return "Internal server error"
			}
		}
	}
	return ""
}

type ResetPasswordForm struct {
	Otp             string `form:"otp"`
	Password        string `form:"password1"`
	ConfirmPassword string `form:"password2"`
}

type ConfirmOtpForm struct {
	Otp   string `form:"otp"`
	Email string `form:"email"`
}

type ForgotPasswordForm struct {
	Email string `form:"email"`
}

type SigninForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type SignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func ResetPasswordApi(c echo.Context, env *types.Env) error {
	resetPasswordForm := new(ResetPasswordForm)
	if err := c.Bind(resetPasswordForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}
	fmt.Println("otp =", resetPasswordForm.Otp)

	user, err := env.Users.Db.GetUserUsingOtp(c.Request().Context(), pgtype.Text{String: resetPasswordForm.Otp, Valid: true})
	if err != nil {
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "Incorrect OTP"})
	}

	if resetPasswordForm.Password != resetPasswordForm.ConfirmPassword {
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "Passwords do not match"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPasswordForm.Password), bcrypt.DefaultCost)

	env.Users.Db.ResetUserPassword(c.Request().Context(), db.ResetUserPasswordParams{
		Password: string(hashedPassword),
		Email:    user.Email,
	})

	c.Response().Header().Set("HX-Redirect", "/signin")
	return c.NoContent(200)
	// return c.Render(200, "signin.html", template.Response{
	// Message: "Password reset successfully",
	// Error:   "",
	// })
}

func ConfirmOtpApi(c echo.Context, env *types.Env) error {
	confirmOtpForm := new(ConfirmOtpForm)
	if err := c.Bind(confirmOtpForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	fmt.Println("otp =", confirmOtpForm.Otp)
	user, err := env.Users.Db.GetUserUsingEmail(c.Request().Context(), confirmOtpForm.Email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("user otp = ", user.Otp.String)
	if user.Otp.String == confirmOtpForm.Otp {
		fmt.Println("ok")
	} else {
		fmt.Println("not ok")
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "Incorrect OTP"})
	}
	return c.Render(200, "confirm-password.html", template.ResetPasswordResponse{Otp: user.Otp.String})
}

func ForgotPasswordApi(c echo.Context, env *types.Env) error {
	forgotPasswordForm := new(ForgotPasswordForm)
	if err := c.Bind(forgotPasswordForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	user, _ := env.Users.Db.GetUserUsingEmail(c.Request().Context(), forgotPasswordForm.Email)
	// if err != nil {
	// return c.Render(200, "errors", template.Response{Message: "notok", Error: getErrorStringFromPgxError(err)})
	// }
	if user.Email == "" {
		return c.Render(200, "forgot-password.html", template.Response{Message: "notok", Error: "User does not exist"})
	}

	id, err := gonanoid.New(12)
	if err != nil {
		return c.Render(200, "forgot-password.html", template.Response{Message: "notok", Error: "Internal server error"})
	}
	err = env.Users.Db.UpdateUserOtp(c.Request().Context(), db.UpdateUserOtpParams{
		Email: forgotPasswordForm.Email,
		Otp:   pgtype.Text{String: id, Valid: true},
	})
	if err != nil {
		return c.Render(200, "forgot-password.html", template.Response{Message: "notok", Error: "Internal server error"})
	}
	return c.Render(200, "confirm-otp.html", template.ForgotPasswordSuccessResponse{Email: forgotPasswordForm.Email})
	// return c.Render(200, "confirm-otp.html", nil)
}

func SignInApi(c echo.Context, env *types.Env) error {
	signinForm := new(SigninForm)
	if err := c.Bind(signinForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	user, err := env.Users.Db.GetUserUsingEmail(c.Request().Context(), signinForm.Email)
	// if err != nil {
	// return c.Render(200, "errors", template.Response{Message: "notok", Error: err.Error()})
	// }

	if user.Email == "" {
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "User does not exist"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signinForm.Password))
	if err != nil {
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "Password does not match"})
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
	return c.NoContent(200)
}

func SignUpApi(c echo.Context, env *types.Env) error {
	signupForm := new(SignupForm)
	if err := c.Bind(signupForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}
	fmt.Println(signupForm.Email, signupForm.Password, signupForm.Name)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupForm.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Render(200, "signup.html", template.Response{Message: "notok", Error: "Internal server error"})
	}
	_, err = env.Users.Db.Create(c.Request().Context(), db.CreateParams{
		Name:     pgtype.Text{String: signupForm.Name, Valid: true},
		Email:    signupForm.Email,
		Password: string(hashedPassword),
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
