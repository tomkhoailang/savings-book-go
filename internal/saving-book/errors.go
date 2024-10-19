package saving_book

var(
	NoCurrentTransaction = "there is no ongoing transaction"
	InsufficientBalance = "insufficient balance"
	MinWithdrawValueError = "doesn't reach the minimum withdraw value"
	TransactionTicketNotPendingStatus = "transaction ticket must be pending to process"
	NotSavingBookOwnerError = "must be an owner to confirm payment"
	TransactionTicketNotFound = "there is no transaction ticket with that payment id"
	CannotDepositError = "the saving book must be in expire to deposit"
)
