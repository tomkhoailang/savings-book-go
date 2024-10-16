package websocket

import "encoding/json"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case broadcastMessage := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- broadcastMessage:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) SendAll(topic string, message interface{}) {
	socketMessage := SocketMessage{
		Type: topic,
		Data: message,
	}
	jsonData,err := json.Marshal(socketMessage)
	if err != nil {
		return
	}
	h.broadcast <- jsonData
	return
}
