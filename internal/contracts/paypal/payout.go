package paypal

type UCPayoutRequest struct {
	Amount float64 `json:"amount" validate:"required,min=1"`
	Email  string  `json:"email" validate:"required,email"`
}
type PaypalTokenResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
	Nonce       string `json:"nonce"`
}
type PayoutRequest struct {
	SenderBatchHeader SenderBatchHeader `json:"sender_batch_header"`
	Items             []PayoutItem      `json:"items"`
}

type PayoutItem struct {
	RecipientType        string       `json:"recipient_type"`
	Amount               PayoutAmount `json:"amount"`
	Note                 string       `json:"note"`
	SenderItemID         string       `json:"sender_item_id"`
	Receiver             string       `json:"receiver"`
	NotificationLanguage string       `json:"notification_language"`
}
type PayoutAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type BatchHeader struct {
	SenderBatchHeader SenderBatchHeader `json:"sender_batch_header"`
	PayoutBatchID     string            `json:"payout_batch_id"`
	BatchStatus       string            `json:"batch_status"`
}

type PayoutBatchResponse struct {
	BatchHeader BatchHeader `json:"batch_header"`
}

type SenderBatchHeader struct {
	SenderBatchID string `json:"sender_batch_id"`
	EmailSubject  string `json:"email_subject"`
	EmailMessage  string `json:"email_message"`
}
