package presenter
type LogInResponse struct {
	Token string `json:"token" `
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
