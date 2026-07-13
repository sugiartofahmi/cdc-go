package dtos

import (
	"time"

	"go-service/entities"

	"github.com/google/uuid"
)

type ProductCategoryResponseDto struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductResponseDto struct {
	Id          uuid.UUID                   `json:"id"`
	Name        string                      `json:"name"`
	Slug        string                      `json:"slug"`
	Description string                      `json:"description"`
	Price       float64                     `json:"price"`
	Stock       int                         `json:"stock"`
	CategoryId  uuid.UUID                   `json:"category_id"`
	Category    *ProductCategoryResponseDto `json:"category,omitempty"`
	CreatedAt   *time.Time                  `json:"created_at"`
	UpdatedAt   *time.Time                  `json:"updated_at"`
}

func ProductResponseDtoFromEntity(entity *entities.ProductEntity) *ProductResponseDto {
	dto := &ProductResponseDto{
		Id:          entity.Id,
		Name:        entity.Name,
		Slug:        entity.Slug,
		Description: entity.Description,
		Price:       entity.Price,
		Stock:       entity.Stock,
		CategoryId:  entity.CategoryId,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}

	if entity.Category != nil {
		dto.Category = &ProductCategoryResponseDto{
			Id:   entity.Category.Id,
			Name: entity.Category.Name,
		}
	}

	return dto
}

func ProductResponseDtoFromEntities(entities []*entities.ProductEntity) []ProductResponseDto {
	responses := make([]ProductResponseDto, len(entities))
	for i, product := range entities {
		responses[i] = *ProductResponseDtoFromEntity(product)
	}
	return responses
}
