package interfaces

import (
	"context"

	"go-service/domain/product/dtos"
)

type ProductEventServiceInterface interface {
	Handle(ctx context.Context, dto *dtos.ProductEventDto)
}
