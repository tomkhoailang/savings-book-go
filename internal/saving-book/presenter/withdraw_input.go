package presenter

type WithDrawInput struct {
	Amount float64 `json:"amount" validate:"required,min=1"`
	Email  string  `json:"email" validate:"required,email"`
}
