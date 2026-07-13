package dtos

import (
	"time"

	"go-service/entities"

	"github.com/google/uuid"
)

type RoleResponseDto struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func RoleResponseDtoFromEntity(role *entities.RoleEntity) *RoleResponseDto {
	return &RoleResponseDto{
		Id:        role.Id,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func RoleResponseDtoFromEntities(entities []*entities.RoleEntity) []RoleResponseDto {
	responses := make([]RoleResponseDto, len(entities))
	for i, role := range entities {
		responses[i] = *RoleResponseDtoFromEntity(role)
	}
	return responses
}
