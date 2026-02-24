package paypal



// ============================================
// Request
// ============================================
type CaptureOrderRequest struct {
	OrderId string `json:"order_id"`
}

type InitOrderRequest struct {
	SavingBookId string `json:"saving_book_id" validate:"required"`
	Amount       string `json:"amount" validate:"required"`
}

type CreateOrderRequest struct {
	Intent        string         `json:"intent"`
	PurchaseUnits []PurchaseUnit `json:"purchase_units"`
	PaymentSource PaymentSource  `json:"payment_source"`
}

type PurchaseUnit struct {
	ReferenceId string `json:"reference_id"`
	Amount      Amount `json:"amount"`
}

type Amount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type ExperienceContext struct {
	PaymentMethodPreference string `json:"payment_method_preference"`
	Locale                  string `json:"locale"`
	ShippingPreference      string `json:"shipping_preference"`
	UserAction              string `json:"user_action"`
	ReturnURL               string `json:"return_url"`
	CancelURL               string `json:"cancel_url"`
}


type PaypalForIntent struct {
	ExperienceContext ExperienceContext `json:"experience_context"`
}
// ============================================
// Response
// ============================================


type PayPalOrderResponse struct {
	Id           string        `json:"id"`
	Status       string        `json:"status"`
	PaymentSource PaymentSource `json:"payment_source"`
	Links        []Link        `json:"links"`
}

type PayPalCaptureResponse struct {
	Id            string         `json:"id"`
	Status        string         `json:"status"`
	PaymentSource PaymentSource  `json:"payment_source"`
	PurchaseUnits []PurchaseUnit `json:"purchase_units"`
	Payer         Payer          `json:"payer"`
	Links         []Link         `json:"links"`
}

type PaymentSource struct {
	PayPal PaypalForIntent `json:"paypal"`
}

type PayPal struct {
	EmailAddress   string `json:"email_address"`
	AccountId      string `json:"account_id"`
	AccountStatus  string `json:"account_status"`
	Name           Name   `json:"name"`
	Address        Address `json:"address"`
}

type Name struct {
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
}

type Address struct {
	CountryCode string `json:"country_code"`
}

type Payments struct {
	Captures []Capture `json:"captures"`
}

type Capture struct {
	Id                     string                  `json:"id"`
	Status                 string                  `json:"status"`
	Amount                 Amount                  `json:"amount"`
	FinalCapture           bool                    `json:"final_capture"`
	SellerProtection       SellerProtection        `json:"seller_protection"`
	SellerReceivableBreakdown SellerReceivableBreakdown `json:"seller_receivable_breakdown"`
	Links                  []Link                  `json:"links"`
	CreateTime             string                  `json:"create_time"`
	UpdateTime             string                  `json:"update_time"`
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

type Payer struct {
	Name         Name    `json:"name"`
	EmailAddress string  `json:"email_address"`
	PayerId      string  `json:"payer_id"`
	Address      Address `json:"address"`
}

