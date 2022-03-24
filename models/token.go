package models

import (
	"github.com/google/uuid"
)

type Token struct {
	UUID   uuid.UUID `json:"id"`
	UserId int       `json:"user_id"`
}
