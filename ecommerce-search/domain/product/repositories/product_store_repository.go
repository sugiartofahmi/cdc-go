package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	productDtos "go-service/domain/product/dtos"
	productInterfaces "go-service/domain/product/interfaces"
	"go-service/infrastructure/exceptions"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

const storeProductIndex = "products-index"

type ProductStoreRepository struct {
	client *opensearchapi.Client
}

func NewProductStoreRepository(client *opensearchapi.Client) productInterfaces.ProductStoreRepositoryInterface {
	return &ProductStoreRepository{client: client}
}

func (r *ProductStoreRepository) Upsert(ctx context.Context, dto *productDtos.ProductResultDto) {
	if r.client == nil {
		log.Println("Error open search client not initialized")
		panic(*exceptions.ServiceUnavailableException("opensearch client not initialized"))
	}

	body, err := json.Marshal(dto)
	if err != nil {
		log.Println("Error marshal product dto:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	req := opensearchapi.IndexReq{
		Index:      storeProductIndex,
		DocumentID: dto.Id,
		Body:       bytes.NewReader(body),
	}

	_, err = r.client.Index(ctx, req)
	if err != nil {
		log.Println("Error upsert product to opensearch:", err)
		panic(*exceptions.ServerErrorException(err))
	}
}

func (r *ProductStoreRepository) Delete(ctx context.Context, id string) {
	if r.client == nil {
		log.Println("Error open search client not initialized")
		panic(*exceptions.ServiceUnavailableException("opensearch client not initialized"))
	}

	req := opensearchapi.DocumentDeleteReq{
		Index:      storeProductIndex,
		DocumentID: id,
	}

	_, err := r.client.Document.Delete(ctx, req)
	if err != nil {
		log.Println("Error delete product from opensearch:", err)
		panic(*exceptions.ServerErrorException(err))
	}
}
