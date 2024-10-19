package presenter

type AuthData struct {
	UserId string `json:"user_id"`
	Roles map[string]interface{} `json:"roles"`
}
