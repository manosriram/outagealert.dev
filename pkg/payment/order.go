package payment

import (
	"encoding/json"
	"fmt"
	"os"

	cashfree "github.com/cashfree/cashfree-pg/v3"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"

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
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Internal server error"})
	}

	free := "free"
	if *user.Plan != free {
		c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "errors", template.Response{Error: "Existing plan active for user"})
	}

	if os.Getenv("ENV") == "production" {
		cashfree.XEnvironment = cashfree.PRODUCTION

		clientId := os.Getenv("CASHFREE_CLIENT_ID")
		clientSecret := os.Getenv("CASHFREE_SECRET_KEY")
		cashfree.XClientId = &clientId
		cashfree.XClientSecret = &clientSecret
	} else {
		cashfree.XEnvironment = cashfree.SANDBOX

		clientId := os.Getenv("CASHFREE_TEST_CLIENT_ID")
		clientSecret := os.Getenv("CASHFREE_TEST_SECRET_KEY")
		cashfree.XClientId = &clientId
		cashfree.XClientSecret = &clientSecret
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
	returnUrl := fmt.Sprintf("%s/projects", os.Getenv("HOST_WITH_SCHEME"))
	request := cashfree.CreateOrderRequest{
		OrderId:     &orderId,
		OrderAmount: amount,
		CustomerDetails: cashfree.CustomerDetails{
			CustomerId:    user.Uuid,
			CustomerEmail: &user.Email,
			CustomerName:  user.Name,
			CustomerPhone: "+917013090094",
		},
		OrderMeta: &cashfree.OrderMeta{
			ReturnUrl: &returnUrl,
		},
		OrderCurrency: "INR",
		OrderSplits:   []cashfree.VendorSplit{},
	}
	version := "2023-08-01"
	orderEntity, httpres, err := cashfree.PGCreateOrder(&version, &request, nil, nil, nil)
	if err != nil {
		l.Log.Error(err.Error())
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
			c.Response().Header().Set("HX-Redirect", "/signin")
			l.Log.Errorf("Error creating new order %s", err.Error())
		}

		// l.Log.Errorf("Error creating cashfree order - %s", err.Error())
		// c.Response().Header().Set("HX-Retarget", "#error-container")
		return c.Render(200, "pricing.html", template.OrderCreatedResponse{PaymentSessionId: *o.PaymentSessionId, ENV: os.Getenv("ENV")})
	} else {
		fmt.Println("order = ", orderEntity)
	}
	return nil
}

type WebhookPayload struct {
	Data      PaymentData `json:"data"`
	EventTime string      `json:"event_time"`
	Type      string      `json:"type"`
}

type PaymentData struct {
	Order                 Order                 `json:"order"`
	Payment               Payment               `json:"payment"`
	CustomerDetails       CustomerDetails       `json:"customer_details"`
	PaymentGatewayDetails PaymentGatewayDetails `json:"payment_gateway_details"`
	PaymentOffers         interface{}           `json:"payment_offers"`
}

type Order struct {
	OrderID       string      `json:"order_id"`
	OrderAmount   float64     `json:"order_amount"`
	OrderCurrency string      `json:"order_currency"`
	OrderTags     interface{} `json:"order_tags"`
}

type Payment struct {
	CFPaymentID     string        `json:"cf_payment_id"`
	PaymentStatus   string        `json:"payment_status"`
	PaymentAmount   float64       `json:"payment_amount"`
	PaymentCurrency string        `json:"payment_currency"`
	PaymentMessage  interface{}   `json:"payment_message"`
	PaymentTime     string        `json:"payment_time"`
	BankReference   interface{}   `json:"bank_reference"`
	AuthID          interface{}   `json:"auth_id"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	PaymentGroup    string        `json:"payment_group"`
}

type PaymentMethod struct {
	UPI UPI `json:"upi"`
}

type UPI struct {
	Channel interface{} `json:"channel"`
	UPIID   interface{} `json:"upi_id"`
}

type CustomerDetails struct {
	CustomerName  string `json:"customer_name"`
	CustomerID    string `json:"customer_id"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

type PaymentGatewayDetails struct {
	GatewayName             string      `json:"gateway_name"`
	GatewayOrderID          string      `json:"gateway_order_id"`
	GatewayPaymentID        string      `json:"gateway_payment_id"`
	GatewayStatusCode       interface{} `json:"gateway_status_code"`
	GatewayOrderReferenceID string      `json:"gateway_order_reference_id"`
	GatewaySettlement       string      `json:"gateway_settlement"`
}

func OrderWebhook(c echo.Context, env *types.Env) error {
	webhookResponse := new(WebhookPayload)
	if err := c.Bind(webhookResponse); err != nil {
		l.Log.Errorf("Error binding payment webhook {}", err.Error())
		return err
	}

	orderId := webhookResponse.Data.Order.OrderID
	order, err := env.DB.Query.GetOrderByOrderId(c.Request().Context(), orderId)
	if err != nil {
		l.Log.Errorf("Error getting order by order_id {}", err.Error())
		return c.JSON(500, nil)
	}

	status := webhookResponse.Data.Payment.PaymentStatus

	// TODO: update to get plan via webhook as well
	var plan string
	switch webhookResponse.Data.Order.OrderAmount {
	case 300:
		plan = "hobbyist"
	case 850:
		plan = "pro"
	}
	if status == "SUCCESS" {
		err = env.DB.Query.UpdateUserPlan(c.Request().Context(), db.UpdateUserPlanParams{
			Plan:  &plan,
			Email: order.UserEmail,
		})
		if err != nil {
			l.Log.Errorf("Error updating user plan {}", err.Error())
			return c.JSON(500, nil)
		}
	}

	j, err := json.Marshal(webhookResponse)
	if err != nil {
		l.Log.Errorf("Error unmarshalling payment webhook {}", err.Error())
		return c.JSON(500, nil)
	}

	err = env.DB.Query.UpdateOrderStatusAndMetadata(c.Request().Context(), db.UpdateOrderStatusAndMetadataParams{
		OrderID:       webhookResponse.Data.Order.OrderID,
		OrderStatus:   webhookResponse.Data.Payment.PaymentStatus,
		OrderMetadata: j,
	})
	if err != nil {
		l.Log.Errorf("Error updating webhook order via query{}", err.Error())
		return c.JSON(500, nil)
	}

	return c.JSON(200, nil)
}
