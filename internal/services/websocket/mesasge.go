package websocket

type SocketMessage struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}