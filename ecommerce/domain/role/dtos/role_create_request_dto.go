package dtos

type RoleCreateRequestDto struct {
	Name string `json:"name" binding:"required"`
}
