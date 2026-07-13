package dtos

import "github.com/google/uuid"

type ProductCreateRequestDto struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gte=0"`
	Stock       int       `json:"stock" binding:"gte=0"`
	CategoryId  uuid.UUID `json:"category_id" binding:"required"`
}
