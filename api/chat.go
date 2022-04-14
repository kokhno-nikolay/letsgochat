package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request, token string) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: ", err.Error())
		return
	}
	defer func() {
		ws.Close()
	}()

	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		log.Println("New message: ", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func (h *Handler) checkToken(token string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.sessions[token]
	return ok
}

func (h *Handler) deleteToken(token string) error {
	if err := h.userRepo.SwitchToInactive(h.sessions[token]); err != nil {
		return err
	}

	_, ok := h.sessions[token]
	if ok {
		h.mu.Lock()
		delete(h.sessions, token)
		h.mu.Unlock()
	}

	return nil
}

func (h *Handler) checkUserSession(userId int) (bool, string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for key, value := range h.sessions {
		if value == userId {
			return true, key
		}
	}

	return false, ""
}
