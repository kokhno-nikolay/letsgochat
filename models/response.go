package models

type JSONResult struct {
	Message string      `json:"message" example:"description of the response from the server"`
	Data    interface{} `json:"data"`
}
