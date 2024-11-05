package presenter

type ResetPasswordConfirm struct {
	Password string `json:"password" validate:"required,min=6,max=20,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}
