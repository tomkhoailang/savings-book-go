package presenter

type RenewTokenReq struct {
	UserId       string `json:"userId" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}
