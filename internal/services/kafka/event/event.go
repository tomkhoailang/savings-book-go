package event

type WithDrawEvent struct {
	Amount float64 `json:"amount"`
	SavingBookId string `json:"savingBookId"`
	TransactionId string `json:"transactionId"`
	Email string `json:"email"`
}
