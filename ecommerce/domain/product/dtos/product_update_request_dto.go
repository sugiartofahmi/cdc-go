package dtos

import "github.com/google/uuid"

type ProductUpdateRequestDto struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Price       *float64   `json:"price" binding:"omitempty,gte=0"`
	Stock       *int       `json:"stock" binding:"omitempty,gte=0"`
	CategoryId  *uuid.UUID `json:"category_id"`
}
