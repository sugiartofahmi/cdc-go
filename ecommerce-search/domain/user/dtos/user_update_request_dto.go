package dtos

import "github.com/google/uuid"

type UserUpdateRequestDto struct {
	Name     *string    `json:"name"`
	Password *string    `json:"password"`
	Status   *int       `json:"status"`
	RoleId   *uuid.UUID `json:"role_id"`
}
