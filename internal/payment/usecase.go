package payment

import (
	"context"

	"SavingBooks/internal/contracts/paypal"
)

type PaymentUseCase interface {
	SendPayout(ctx context.Context, payoutRequest *paypal.UCPayoutRequest) (*paypal.PayoutBatchResponse, error)
	CreateOrder(ctx context.Context, orderRequest *paypal.InitOrderRequest) (*paypal.PayPalOrderResponse, error)
	CaptureOrder(ctx context.Context, orderId string) (*paypal.PayPalCaptureResponse, error)
}