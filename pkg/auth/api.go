package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/integration"
	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

const (
	NANOID_ALPHABET_LIST = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"
	NANOID_LENGTH        = 22
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
	host := os.Getenv("HOST_WITH_SCHEME")
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("%s/signin", host))

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
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Incorrect OTP"})
	}

	if resetPasswordForm.Password != resetPasswordForm.ConfirmPassword {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Passwords do not match"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPasswordForm.Password), bcrypt.DefaultCost)

	env.DB.Query.ResetUserPassword(c.Request().Context(), db.ResetUserPasswordParams{
		Password: string(hashedPassword),
		Email:    user.Email,
	})

	err = c.Render(200, "signin.html", template.Response{
		Message: "Password reset successfully",
	})
	if err != nil {
		l.Log.Error(err)
	}
	return err
}

func ConfirmOtpApi(c echo.Context, env *types.Env) error {
	confirmOtpForm := new(ConfirmOtpForm)
	if err := c.Bind(confirmOtpForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), confirmOtpForm.Email)
	if err != nil {
		l.Log.Error(err)
		return err
	}
	if *user.Otp != confirmOtpForm.Otp {
		return c.Render(200, "errors", template.Response{Error: "Incorrect OTP"})
	}

	err = env.DB.Query.UpdateUserMagicToken(c.Request().Context(), db.UpdateUserMagicTokenParams{
		MagicToken: &confirmOtpForm.Otp,
		ID:         user.ID,
	})
	if err != nil {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Error registering user"})
	}

	return c.Render(200, "confirm-password.html", template.ResetPasswordResponse{Otp: *user.Otp})
}

func ForgotPasswordApi(c echo.Context, env *types.Env) error {
	forgotPasswordForm := new(ForgotPasswordForm)
	if err := c.Bind(forgotPasswordForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	validations := env.Validator.Struct(forgotPasswordForm)
	if validations != nil && len(validations.(validator.ValidationErrors)) > 0 {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		switch validations.(validator.ValidationErrors)[0].Tag() {
		case "email":
			return c.Render(200, "errors", template.Response{Error: "Invalid email"})
		default:
			return c.Render(200, "errors", template.Response{Error: "Invalid details"})
		}
	}

	// todo: handle err
	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), forgotPasswordForm.Email)
	if err != nil {
		l.Log.Errorf("Error getting user email %s", err.Error())
		return c.Render(200, "errors", template.Response{Error: "User does not exist"})
	}

	if user.Email == "" {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "User does not exist"})
	}

	id, err := gonanoid.New(12)
	if err != nil {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}
	err = env.DB.Query.UpdateUserOtp(c.Request().Context(), db.UpdateUserOtpParams{
		Email: forgotPasswordForm.Email,
		Otp:   &id,
	})
	if err != nil {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "User does not exist"})
	}

	notif := integration.EmailNotification{
		Email: forgotPasswordForm.Email,
		OTP:   id,
	}
	go notif.SendMail("forgot_password_otp", "d-038cf4d4bd6a492ca28d19f6d8fe3b24", integration.VerifyEmailMailData{
		OTP: id,
	})

	return c.Render(200, "confirm-otp.html", template.ForgotPasswordSuccessResponse{Email: forgotPasswordForm.Email})
}

func SignInApi(c echo.Context, env *types.Env) error {
	signinForm := new(SigninForm)
	if err := c.Bind(signinForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid form data")
	}

	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), signinForm.Email)
	if err != nil {
		l.Log.Error(err.Error())
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "User not found"})
	}

	if !user.IsVerified {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "User email not verified"})
	}

	if user.Email == "" {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "User does not exist"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signinForm.Password))
	if err != nil {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Password does not match"})
	}

	sess, err := session.Get("session", c)
	if err != nil {
		l.Log.Error(err)
		c.Response().Header().Set("HX-Retarget", "#error-container")
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
		l.Log.Error(err)
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return err
	}
	return c.NoContent(200)
}

func SignUpApi(c echo.Context, env *types.Env) error {
	signupForm := new(SignupForm)
	if err := c.Bind(signupForm); err != nil {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Invalid form data"})
	}

	validations := env.Validator.Struct(signupForm)
	if validations != nil && len(validations.(validator.ValidationErrors)) > 0 {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		switch validations.(validator.ValidationErrors)[0].Tag() {
		case "email":
			return c.Render(200, "errors", template.Response{Error: "Invalid email"})
		case "name":
			return c.Render(200, "errors", template.Response{Error: "Invalid name"})
		case "password":
			return c.Render(200, "errors", template.Response{Error: "Invalid password"})
		default:
			return c.Render(200, "errors", template.Response{Error: "Invalid details"})
		}
	}

	user, _ := env.DB.Query.GetUserUsingEmail(c.Request().Context(), signupForm.Email)
	if user.Email != "" {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "User already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupForm.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	magicToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": signupForm.Email,
		"exp":   time.Now().Add(time.Minute * 10).Unix(),
	})
	tokenString, err := magicToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	encToken := base64.StdEncoding.EncodeToString([]byte(tokenString))
	notif := integration.EmailNotification{
		Email:      signupForm.Email,
		MagicToken: encToken,
	}

	createdUser, err := env.DB.Query.Create(c.Request().Context(), db.CreateParams{
		Name:       &signupForm.Name,
		Email:      signupForm.Email,
		Password:   string(hashedPassword),
		MagicToken: &encToken,
		Uuid:       gonanoid.Must(NANOID_LENGTH),
	})
	if err != nil {
		l.Log.Errorf("Error registering user %s", err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				c.Response().Header().Set("HX-Retarget", "#error-container")
				return c.Render(200, "errors", template.Response{Error: "User already exists"})
			default:
				c.Response().Header().Set("HX-Retarget", "#error-container")
				return c.Render(200, "signup.html", template.Response{Error: "Internal server error"})
			}
		}
	}
	l.Log.Infof("User created %s", createdUser.Email)

	// TODO: use interfaced struct to organize different email senders
	go notif.SendMail("verify_email", "d-c50ac0a5dccb454fbbb6eac650b5e680", integration.VerifyEmailMailData{
		Name:      signupForm.Name,
		Subject:   "Reset password - outagealert",
		MagicLink: encToken,
	})
	return c.Render(200, "signup-success", template.RegisterSuccessResponse{Email: signupForm.Email})
}

func VerifyEmailViaMagicToken(c echo.Context, env *types.Env) error {
	magictoken := c.Param("magic_token")
	decodedString, _ := base64.StdEncoding.DecodeString(magictoken)

	token, err := jwt.Parse(string(decodedString), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		l.Log.Error(err)
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email := claims["email"].(string)
		err := env.DB.Query.MarkUserVerified(c.Request().Context(), email)
		if err != nil {
			l.Log.Error(err)
			return c.Render(200, "errors", template.Response{Error: "Internal server error"})
		}
	} else {
		l.Log.Error(err)
	}

	err = c.Redirect(301, fmt.Sprintf("%s/email-verified", os.Getenv("HOST_WITH_SCHEME")))
	if err != nil {
		l.Log.Error(err)
	}
	return err
}
