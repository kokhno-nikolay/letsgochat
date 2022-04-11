package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
