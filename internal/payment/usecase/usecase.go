package usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	paypal "SavingBooks/internal/contracts/paypal"
	"SavingBooks/internal/payment"
	"github.com/google/uuid"
)

type paymentUseCase struct {
	clientId string
	clientSecret string
	httpClient *http.Client
}

func (p *paymentUseCase) SendPayout(ctx context.Context, payoutRequest *paypal.UCPayoutRequest) (*paypal.PayoutBatchResponse, error) {
	token, err := p.getPaypalToken()
	if err != nil {
		return nil, err
	}
	payload := &paypal.PayoutRequest{
		SenderBatchHeader: paypal.SenderBatchHeader{
			SenderBatchID: fmt.Sprintf("PayoutsBatch_%s_%s",time.Now(),uuid.New().String()),
			EmailSubject: "You have a payout",
			EmailMessage: "You have received a payout! Thanks for using our service!",
		},
		Items: []paypal.PayoutItem {
			{
				RecipientType: "EMAIL",
				Amount: paypal.PayoutAmount{
					Value: fmt.Sprintf("%.2f", math.Ceil(payoutRequest.Amount*100)/100), Currency: "USD" ,
				},
				Note: "Thanks for your patronage!",
				SenderItemID: fmt.Sprintf("Batch_%s",time.Now().Format("20060102 150405")),
				Receiver: payoutRequest.Email,
				NotificationLanguage: "en-US",
			},
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api-m.sandbox.paypal.com/v1/payments/payouts", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil,err
		}
		return nil, errors.New(string(body))
	}
	defer resp.Body.Close()
	var payoutBatchRes paypal.PayoutBatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payoutBatchRes); err != nil {
		return nil, err
	}
	return &payoutBatchRes, nil
}

func (p *paymentUseCase) CreateOrder(ctx context.Context, orderRequest *paypal.InitOrderRequest) (*paypal.PayPalOrderResponse, error) {
	token, err := p.getPaypalToken()
	if err != nil {
		return nil, err
	}
	payload := &paypal.CreateOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []paypal.PurchaseUnit{
			{
				ReferenceId: orderRequest.SavingBookId,
				Amount: paypal.Amount{CurrencyCode: "USD", Value: orderRequest.Amount},
			},
		},
		PaymentSource: paypal.PaymentSource{
			PayPal: paypal.PaypalForIntent{
				ExperienceContext: paypal.ExperienceContext{
					PaymentMethodPreference: "IMMEDIATE_PAYMENT_REQUIRED",
					Locale:                  "en-US",
					ShippingPreference:      "NO_SHIPPING",
					UserAction:              "PAY_NOW",
					ReturnURL:               "https://example.com/returnUrl",
					CancelURL:               "https://example.com/cancelUrl",
				},
			},
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api-m.sandbox.paypal.com/v2/checkout/orders", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil,err
		}
		return nil, errors.New(string(body))
	}
	defer resp.Body.Close()
	var paypalResp paypal.PayPalOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&paypalResp); err != nil {
		return nil, err
	}
	return &paypalResp, nil
}

func (p *paymentUseCase) CaptureOrder(ctx context.Context, orderId string) (*paypal.PayPalCaptureResponse, error) {
	token, err := p.getPaypalToken()
	if err != nil {
		return nil, err
	}

	jsonPayload, err := json.Marshal(struct{}{})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api-m.sandbox.paypal.com/v2/checkout/orders/%s/capture", orderId), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil,err
		}
		return nil, errors.New(string(body))
	}
	defer resp.Body.Close()
	var paypalResp paypal.PayPalCaptureResponse
	if err := json.NewDecoder(resp.Body).Decode(&paypalResp); err != nil {
		return nil, err
	}
	return &paypalResp, nil
}

func (p *paymentUseCase) getPaypalToken() (string, error) {
	auth := p.clientId + ":" + p.clientSecret
	authToken := base64.StdEncoding.EncodeToString([]byte(auth))
	req, err := http.NewRequest(http.MethodPost, "https://api-m.sandbox.paypal.com/v1/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error getting PayPal access token: %s", string(body))
	}
	var tokenResult paypal.PaypalTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResult); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResult.AccessToken, nil

}

func NewPaymentUseCase(clientId string, clientSecret string) payment.PaymentUseCase {
	return &paymentUseCase{
		clientId: clientId,
		clientSecret: clientSecret,
		httpClient: &http.Client{},
	}
}
