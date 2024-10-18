package presenter

type DepositInput struct {
	Amount float64 `json:"amount" validate:"required,min=1"`
}