package api

import (
	"io"
	"log"
	"net/http"

	"github.com/kokhno-nikolay/letsgochat/models"

	"github.com/gorilla/websocket"
)

type message struct {
	Client  *websocket.Conn
	Message models.ChatMessage
}

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

		if len(msg.Text) < 1 {
			log.Println("incorrect message text")
			continue
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

func (h *Handler) workerMessage(msg models.ChatMessage) {
	h.messageCh = make(chan message, workersNum)

	for client := range h.clients {
		msg := message{Client: client, Message: msg}
		h.messageCh <- msg
	}
	close(h.messageCh)
}

func (h *Handler) messageClients() {
	for msg := range h.messageCh {
		h.messageClient(msg.Client, msg.Message)
	}
}

func (h *Handler) unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func (h *Handler) handleMessages(token string) {
	for {
		msg := <-h.broadcaster

		user, err := h.userRepo.FindById(h.Sessions[token])
		if err != nil {
			log.Fatal(err.Error())
		}

		msgModel := models.Message{Text: msg.Text, UserId: user.ID}
		if err := h.messageRepo.Create(msgModel); err != nil {
			log.Fatal(err.Error())
		}

		msg.Username = user.Username
		h.workerMessage(msg)

		if len(h.clients) < workersNum {
			for i := 0; i <= len(h.clients); i++ {
				go h.messageClients()
			}
		} else {
			for i := 0; i <= workersNum; i++ {
				go h.messageClients()
			}
		}
	}
}
