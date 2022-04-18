package api

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/kokhno-nikolay/letsgochat/models"
)

func (h *Handler) handleConnections(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	h.clients[ws] = true

	messages, err := h.messageRepo.GetAll()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(messages) != 0 {
		for _, msg := range messages {
			h.messageClient(ws, msg)
		}
	}

	for {
		var msg models.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(h.clients, ws)
			break
		}

		h.broadcaster <- msg
	}
}

func (h *Handler) messageClient(client *websocket.Conn, msg models.ChatMessage) {
	err := client.WriteJSON(msg)
	if err != nil && h.unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()
		delete(h.clients, client)
	}
}

func (h *Handler) messageClients(msg models.ChatMessage) {
	for client := range h.clients {
		h.messageClient(client, msg)
	}
}

func (h *Handler) unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func (h *Handler) handleMessages(token string) {
	for {
		msg := <-h.broadcaster

		user, err := h.userRepo.FindById(h.sessions[token])
		if err != nil {
			log.Fatal(err.Error())
		}

		msgModel := models.Message{Text: msg.Text, UserId: user.ID}
		if err := h.messageRepo.Create(msgModel); err != nil {
			log.Fatal(err.Error())
		}

		msg.Username = user.Username
		h.messageClients(msg)
	}
}
