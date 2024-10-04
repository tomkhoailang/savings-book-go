package presenter

type SignUpInput struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignUpResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}