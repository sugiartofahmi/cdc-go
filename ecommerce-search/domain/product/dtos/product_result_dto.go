package dtos

import (
	"encoding/json"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type ProductResultDto struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CategorySlug string  `json:"category_slug"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func ProductResultDtoFromDocument(source json.RawMessage) *ProductResultDto {
	var dto ProductResultDto
	if err := json.Unmarshal(source, &dto); err != nil {
		return nil
	}
	return &dto
}

func ProductResultDtoFromDocuments(hits []opensearchapi.SearchHit) []*ProductResultDto {
	result := make([]*ProductResultDto, 0, len(hits))
	for _, hit := range hits {
		if dto := ProductResultDtoFromDocument(hit.Source); dto != nil {
			result = append(result, dto)
		}
	}
	return result
}
