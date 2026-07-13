package dtos

import (
	"time"

	"go-service/entities"

	"github.com/google/uuid"
)

type CategoryResponseDto struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func CategoryResponseDtoFromEntity(entity *entities.CategoryEntity) *CategoryResponseDto {
	return &CategoryResponseDto{
		Id:        entity.Id,
		Name:      entity.Name,
		Slug:      entity.Slug,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func CategoryResponseDtoFromEntities(entities []*entities.CategoryEntity) []CategoryResponseDto {
	responses := make([]CategoryResponseDto, len(entities))
	for i, cat := range entities {
		responses[i] = *CategoryResponseDtoFromEntity(cat)
	}
	return responses
}
