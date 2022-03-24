package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
}

type UserInput struct {
	Username string `json:"username" binding:"required,min=4,max=64"`
	Password string `json:"password"  binding:"required,min=8,max=64"`
}
