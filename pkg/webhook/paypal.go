package webhook

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/template"
)

type PaypalWebhookEvent struct {
	ID              string    `json:"id"`
	CreateTime      time.Time `json:"create_time"`
	EventType       string    `json:"event_type"`
	ResourceType    string    `json:"resource_type"`
	ResourceVersion string    `json:"resource_version"`
	Summary         string    `json:"summary"`
	Resource        Resource  `json:"resource"`
	Links           []Link    `json:"links"`
	Zts             int64     `json:"zts"`
	EventVersion    string    `json:"event_version"`
}

type Resource struct {
	ID            string         `json:"id"`
	Status        string         `json:"status"`
	Intent        string         `json:"intent"`
	GrossAmount   Amount         `json:"gross_amount"`
	Payer         Payer          `json:"payer"`
	PurchaseUnits []PurchaseUnit `json:"purchase_units"`
	CreateTime    time.Time      `json:"create_time"`
	UpdateTime    time.Time      `json:"update_time"`
	Links         []Link         `json:"links"`
}

type Amount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type Payer struct {
	Name         Name   `json:"name"`
	EmailAddress string `json:"email_address"`
	PayerID      string `json:"payer_id"`
}

type Name struct {
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
}

type PurchaseUnit struct {
	ReferenceID string   `json:"reference_id"`
	Amount      Amount   `json:"amount"`
	Payee       Payee    `json:"payee"`
	Shipping    Shipping `json:"shipping"`
	Payments    Payments `json:"payments"`
}

type Payee struct {
	EmailAddress string `json:"email_address"`
}

type Shipping struct {
	Method  string  `json:"method"`
	Address Address `json:"address"`
}

type Address struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	AdminArea2   string `json:"admin_area_2"`
	AdminArea1   string `json:"admin_area_1"`
	PostalCode   string `json:"postal_code"`
	CountryCode  string `json:"country_code"`
}

type Payments struct {
	Captures []Capture `json:"captures"`
}

type Capture struct {
	ID                        string                    `json:"id"`
	Status                    string                    `json:"status"`
	Amount                    Amount                    `json:"amount"`
	SellerProtection          SellerProtection          `json:"seller_protection"`
	FinalCapture              bool                      `json:"final_capture"`
	SellerReceivableBreakdown SellerReceivableBreakdown `json:"seller_receivable_breakdown"`
	CreateTime                time.Time                 `json:"create_time"`
	UpdateTime                time.Time                 `json:"update_time"`
	Links                     []Link                    `json:"links"`
}

type SellerProtection struct {
	Status            string   `json:"status"`
	DisputeCategories []string `json:"dispute_categories"`
}

type SellerReceivableBreakdown struct {
	GrossAmount Amount `json:"gross_amount"`
	PaypalFee   Amount `json:"paypal_fee"`
	NetAmount   Amount `json:"net_amount"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

func PaypalWebhook(c echo.Context) error {
	paypalWebhookEvent := new(PaypalWebhookEvent)
	if err := c.Bind(paypalWebhookEvent); err != nil {
		return c.JSON(400, template.Response{Message: "", Error: err.Error()})
	}
	fmt.Println(paypalWebhookEvent.EventType)
	switch paypalWebhookEvent.EventType {
	case "PAYMENT.CAPTURE.COMPLETED":
		break
	case "CHECKOUT.ORDER.APPROVED":
		fmt.Printf("checkout done for %s\n", paypalWebhookEvent.Resource.GrossAmount.Value)
	}

	return c.JSON(200, template.Response{Message: "Paypal event consumed", Error: ""})
}
