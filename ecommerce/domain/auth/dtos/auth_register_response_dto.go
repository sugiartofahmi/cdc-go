package dtos

import (
	"time"

	"github.com/google/uuid"
)

type AuthRegisterResponseDto struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
}
