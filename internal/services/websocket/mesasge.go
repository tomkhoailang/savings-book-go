package websocket

type SocketMessage struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}
type ClientMessage struct {
	ClientID string
	Message  []byte
}
const (
	WithDrawStatus = "WithDrawStatus"
)