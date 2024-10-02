package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgconn"
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
	Email string `form:"email" validate:"required,email"`
}

type ForgotPasswordForm struct {
	Email string `form:"email" validate:"required,email"`
}

type SigninForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password"`
}

type SignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"gte=4"`
}

func Logout(c echo.Context, env *types.Env) error {
	s, err := session.Get("session", c)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	s.Values["email"] = ""
	s.Values["id"] = ""
	s.Options.MaxAge = -1
	s.Save(c.Request(), c.Response())
	c.Response().Header().Set("HX-Redirect", "http://localhost:1323/signin")

	return c.NoContent(200)
}

func GetUser(c echo.Context) (*sessions.Session, error) {
	s, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GetCurrentUser(c echo.Context, env *types.Env) error {
	s, err := session.Get("session", c)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}
	email := s.Values["email"]
	if email == nil {
		return c.NoContent(200)
	}
	fmt.Println(email)

	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), email.(string))
	if err != nil {
		return c.Render(200, "signin", template.User{User: user})
	}

	return c.Render(200, "user-profile", template.User{User: user})
}

func ResetPasswordApi(c echo.Context, env *types.Env) error {
	resetPasswordForm := new(ResetPasswordForm)
	if err := c.Bind(resetPasswordForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	user, err := env.DB.Query.GetUserUsingOtp(c.Request().Context(), &resetPasswordForm.Otp)
	if err != nil {
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "Incorrect OTP"})
	}

	if resetPasswordForm.Password != resetPasswordForm.ConfirmPassword {
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "Passwords do not match"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPasswordForm.Password), bcrypt.DefaultCost)

	env.DB.Query.ResetUserPassword(c.Request().Context(), db.ResetUserPasswordParams{
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

	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), confirmOtpForm.Email)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if *user.Otp != confirmOtpForm.Otp {
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "Incorrect OTP"})
	}
	return c.Render(200, "confirm-password.html", template.ResetPasswordResponse{Otp: *user.Otp})
}

func ForgotPasswordApi(c echo.Context, env *types.Env) error {
	forgotPasswordForm := new(ForgotPasswordForm)
	if err := c.Bind(forgotPasswordForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	var errorMsg string
	validations := env.Validator.Struct(forgotPasswordForm)
	if validations != nil && len(validations.(validator.ValidationErrors)) > 0 {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		switch validations.(validator.ValidationErrors)[0].Tag() {
		case "email":
			errorMsg = "Invalid email"
		default:
			errorMsg = "Invalid details"
		}
		return c.Render(200, "errors", template.Response{Message: "notok", Error: errorMsg})
	}

	// todo: handle err
	user, _ := env.DB.Query.GetUserUsingEmail(c.Request().Context(), forgotPasswordForm.Email)

	if user.Email == "" {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		// return c.Render(200, "errors", template.Response{Message: "notok", Error: "User does not exist"})
		return c.Render(200, "confirm-otp.html", template.ForgotPasswordSuccessResponse{Response: template.Response{Error: "User does not exist"}})
	}

	id, err := gonanoid.New(12)
	if err != nil {
		return c.Render(200, "forgot-password.html", template.Response{Message: "notok", Error: "Internal server error"})
	}
	err = env.DB.Query.UpdateUserOtp(c.Request().Context(), db.UpdateUserOtpParams{
		Email: forgotPasswordForm.Email,
		Otp:   &id,
	})
	if err != nil {
		// return c.Render(200, "forgot-password.html", template.Response{Message: "notok", Error: "Internal server error"})
		return c.Render(200, "forgot-password.html", template.ForgotPasswordSuccessResponse{Response: template.Response{Error: "User does not exist"}})
	}
	return c.Render(200, "confirm-otp.html", template.ForgotPasswordSuccessResponse{Email: forgotPasswordForm.Email})
	// return c.Render(200, "confirm-otp.html", nil)
}

func SignInApi(c echo.Context, env *types.Env) error {
	signinForm := new(SigninForm)
	if err := c.Bind(signinForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), signinForm.Email)
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
	sess.Values["id"] = user.ID
	c.Response().Header().Set("HX-Redirect", "/projects")
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.NoContent(200)
}

func SignUpApi(c echo.Context, env *types.Env) error {
	signupForm := new(SignupForm)
	if err := c.Bind(signupForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	var errorMsg string
	validations := env.Validator.Struct(signupForm)
	if validations != nil && len(validations.(validator.ValidationErrors)) > 0 {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		switch validations.(validator.ValidationErrors)[0].Tag() {
		case "email":
			errorMsg = "Invalid email"
		case "name":
			errorMsg = "Invalid name"
		case "password":
			errorMsg = "Invalid password"
		default:
			errorMsg = "Invalid details"
		}
		return c.Render(200, "errors", template.Response{Message: "notok", Error: errorMsg})
	}

	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), signupForm.Email)
	if user.Email != "" {
		return c.Render(200, "errors", template.Response{Message: "notok", Error: "User already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupForm.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Render(200, "signup.html", template.Response{Message: "notok", Error: "Internal server error"})
	}
	_, err = env.DB.Query.Create(c.Request().Context(), db.CreateParams{
		Name:     &signupForm.Name,
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
	return c.Render(200, "signup-success", template.RegisterSuccessResponse{Email: signupForm.Email})
}
