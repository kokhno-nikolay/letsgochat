package models

type Message struct {
	Id     int    `json:"id"`
	Text   string `json:"text"`
	UserId int    `json:"user_id"`
}

type ChatMessage struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}
