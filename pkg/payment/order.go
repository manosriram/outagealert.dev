package payment

import (
	"encoding/json"
	"fmt"
	"os"

	cashfree "github.com/cashfree/cashfree-pg/v3"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/pkg/template"
	"github.com/manosriram/outagealert.io/pkg/types"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

const (
	NANOID_ALPHABET_LIST = "abcdefghijklmnopqstuvwxyzABCDEFGHIJKLMNOPQSTUVWXYZ"
	NANOID_LENGTH        = 22
)

type PaymentSelect struct {
	Plan string `form:"plan"`
}

func CreateOrder(c echo.Context, env *types.Env) error {
	s, err := session.Get("session", c)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}
	email := s.Values["email"].(string)
	if s.Values["email"] == nil {
		c.Response().Header().Set("HX-Redirect", "/signin")
		return c.NoContent(200)
	}

	plan := c.QueryParam("plan")

	user, err := env.DB.Query.GetUserUsingEmail(c.Request().Context(), email)
	if err != nil {
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	clientId := os.Getenv("CASHFREE_CLIENT_ID")
	clientSecret := os.Getenv("CASHFREE_SECRET_KEY")
	fmt.Println(clientId, clientSecret)
	cashfree.XClientId = &clientId
	cashfree.XClientSecret = &clientSecret
	if os.Getenv("ENV") == "production" {
		cashfree.XEnvironment = cashfree.PRODUCTION
	} else {
		cashfree.XEnvironment = cashfree.SANDBOX
	}
	// TODO: use this via config/env
	var amount float64
	switch plan {
	case "hobbyist":
		amount = 300.00
	case "pro":
		amount = 850.00
	}
	orderId, err := gonanoid.Generate(NANOID_ALPHABET_LIST, NANOID_LENGTH)
	orderId = fmt.Sprintf("%s", orderId)
	request := cashfree.CreateOrderRequest{
		OrderId:     &orderId,
		OrderAmount: amount,
		CustomerDetails: cashfree.CustomerDetails{
			CustomerId:    string(user.ID),
			CustomerEmail: &user.Email,
			CustomerName:  user.Name,
			CustomerPhone: "+917013090094",
		},
		OrderCurrency: "INR",
		OrderSplits:   []cashfree.VendorSplit{},
	}
	version := "2023-08-01"
	orderEntity, httpres, err := cashfree.PGCreateOrder(&version, &request, nil, nil, nil)
	if err != nil {
		fmt.Println(err)
		o := cashfree.OrderEntity{}
		json.NewDecoder(httpres.Body).Decode(&o)

		err = env.DB.Query.CreateOrder(c.Request().Context(), db.CreateOrderParams{
			OrderID:               *o.OrderId,
			UserEmail:             email,
			OrderStatus:           *o.OrderStatus,
			OrderPaymentSessionID: o.PaymentSessionId,
			Plan:                  &plan,
			OrderExpiryTime:       pgtype.Timestamp{Time: *o.OrderExpiryTime, Valid: true},
			OrderCurrency:         o.OrderCurrency,
		})
		if err != nil {
			l.Log.Errorf("Error creating new order %s", err.Error())

		}

		// l.Log.Errorf("Error creating cashfree order - %s", err.Error())
		// c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "pricing.html", template.OrderCreatedResponse{PaymentSessionId: *o.PaymentSessionId})
	} else {
		fmt.Println("order = ", orderEntity)
	}
	return nil
}

type WebhookResponse struct {
	Data      WebhookData `json:"data"`
	Type      string      `json:"type"`
	EventTime string      `json:"event_time"`
}

type WebhookData struct {
	Form  FormData  `json:"form"`
	Order OrderData `json:"order"`
}

type FormData struct {
	CFFormID     int    `json:"cf_form_id"`
	FormCurrency string `json:"form_currency"`
	FormID       string `json:"form_id"`
	FormURL      string `json:"form_url"`
}

type OrderData struct {
	AmountDetails   []AmountDetail  `json:"amount_details"`
	CustomerDetails CustomerDetails `json:"customer_details"`
	OrderAmount     float64         `json:"order_amount"`
	OrderID         string          `json:"order_id"`
	OrderStatus     string          `json:"order_status"`
	TransactionID   int64           `json:"transaction_id"`
}

type AmountDetail struct {
	Quantity       int     `json:"quantity,omitempty"`
	Title          string  `json:"title"`
	Value          float64 `json:"value"`
	SelectedOption string  `json:"selectedoption,omitempty"`
}

type CustomerDetails struct {
	CustomerEmail  string          `json:"customer_email"`
	CustomerFields []CustomerField `json:"customer_fields"`
	CustomerName   string          `json:"customer_name"`
	CustomerPhone  string          `json:"customer_phone"`
}

type CustomerField struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

func OrderWebhook(c echo.Context, env *types.Env) error {
	webhookResponse := new(WebhookResponse)
	if err := c.Bind(webhookResponse); err != nil {
		return err
	}
	fmt.Println("status = ", webhookResponse.Data.Order.OrderStatus)

	return c.Render(200, "errors", template.Response{Error: "Internal server error"})
}
