package websocket

import "encoding/json"

type Hub struct {
	clients    map[string]map[string]*Client
	broadcast  chan []byte
	notify     chan ClientMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		notify:     make(chan ClientMessage),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.userId] == nil {
				h.clients[client.userId] = make(map[string]*Client)
			}
			h.clients[client.userId][client.tabId] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.userId]; ok {
				delete(h.clients, client.userId)
			}
		case broadcastMessage := <-h.broadcast:
			for _, userClients := range h.clients {
				for _, client := range userClients {
					select {
					case client.send <- broadcastMessage:
					default:
						close(client.send)
						delete(h.clients, client.userId)
					}
				}
			}
		case clientMessage := <-h.notify:
			if clients, ok := h.clients[clientMessage.ClientID]; ok {
				for _, client := range clients {
					select {
					case client.send <- clientMessage.Message:
					default:
						close(client.send)
						delete(h.clients, client.userId)
					}
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
	jsonData, err := json.Marshal(socketMessage)
	if err != nil {
		return
	}
	h.broadcast <- jsonData
}
func (h *Hub) SendOne(topic string, clientId string, message interface{}) {
	socketMessage := SocketMessage{
		Type: topic,
		Data: message,
	}
	jsonData, err := json.Marshal(socketMessage)
	if err != nil {
		return
	}
	h.notify <- ClientMessage{
		ClientID: clientId,
		Message:  jsonData,
	}
}
