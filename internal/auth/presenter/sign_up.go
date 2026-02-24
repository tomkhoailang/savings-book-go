package presenter

type SignUpInput struct {
	Username string `json:"username" validate:"required,min=6,max=20"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

type SignUpResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}