package dtos

type CategoryCreateRequestDto struct {
	Name string `json:"name" binding:"required"`
}
