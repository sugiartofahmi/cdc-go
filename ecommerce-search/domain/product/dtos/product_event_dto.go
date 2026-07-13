package dtos

import (
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type ProductEventDto struct {
	Operation string `json:"__op"`

	Id string `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryId  string  `json:"category_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string  `json:"deleted_at"`
}

func ProductEventDtoFromMessage(msg kafka.Message) *ProductEventDto {
	var dto ProductEventDto
	if err := json.Unmarshal(msg.Value, &dto); err != nil {
		return nil
	}
	return &dto
}
