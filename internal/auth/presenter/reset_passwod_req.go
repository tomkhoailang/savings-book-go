package presenter
type ResetPasswordRequest struct  {
	Email string `json:"email" validate:"required,email"`
}
