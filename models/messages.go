package models

type Message struct {
	Id     int    `json:"id"`
	Text   string `json:"text"`
	UserId int    `json:"user_id"`
}
