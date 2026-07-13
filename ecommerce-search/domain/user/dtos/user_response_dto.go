package dtos

import (
	"time"

	"go-service/entities"

	"github.com/google/uuid"
)

type UserRoleResponseDto struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UserResponseDto struct {
	Id        uuid.UUID             `json:"id"`
	Name      string                `json:"name"`
	Email     string                `json:"email"`
	Status    int                   `json:"status"`
	Role      *UserRoleResponseDto  `json:"role"`
	CreatedAt *time.Time            `json:"created_at"`
	UpdatedAt *time.Time            `json:"updated_at"`
}

func UserResponseDtoFromEntity(entity *entities.UserEntity) *UserResponseDto {
	var role *UserRoleResponseDto
	if entity.Role != nil {
		role = &UserRoleResponseDto{
			Id:   entity.Role.Id,
			Name: entity.Role.Name,
		}
	}

	return &UserResponseDto{
		Id:        entity.Id,
		Name:      entity.Name,
		Email:     entity.Email,
		Status:    entity.Status,
		Role:      role,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func UserResponseDtoFromEntities(entities []*entities.UserEntity) []UserResponseDto {
	responses := make([]UserResponseDto, len(entities))
	for i, user := range entities {
		responses[i] = *UserResponseDtoFromEntity(user)
	}
	return responses
}
