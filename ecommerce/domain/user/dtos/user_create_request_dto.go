package dtos

import "github.com/google/uuid"

type UserCreateRequestDto struct {
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	Password string    `json:"password" binding:"required"`
	RoleId   uuid.UUID `json:"role_id" binding:"required"`
}
