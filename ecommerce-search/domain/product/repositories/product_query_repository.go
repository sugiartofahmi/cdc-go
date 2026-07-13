package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	productDtos "go-service/domain/product/dtos"
	productInterfaces "go-service/domain/product/interfaces"
	"go-service/infrastructure/exceptions"
	infradtos "go-service/infrastructure/dtos"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

const productIndex = "products-index"

type ProductQueryRepository struct {
	client *opensearchapi.Client
}

func NewProductQueryRepository(client *opensearchapi.Client) productInterfaces.ProductQueryRepositoryInterface {
	return &ProductQueryRepository{client: client}
}

func (r *ProductQueryRepository) Pagination(ctx context.Context, dto *productDtos.ProductQueryRequestDto) *infradtos.PaginationResultDto[productDtos.ProductResultDto] {
	if r.client == nil {
		log.Println("Error open search client not initialized")
		panic(*exceptions.ServiceUnavailableException("opensearch client not initialized"))
	}

	body, err := r.buildQuery(dto)
	if err != nil {
		log.Println("Error build open search query:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	searchResp, err := r.client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{productIndex},
		Body:    body,
	})
	if err != nil {
		log.Println("Error search products:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	items := productDtos.ProductResultDtoFromDocuments(searchResp.Hits.Hits)

	return &infradtos.PaginationResultDto[productDtos.ProductResultDto]{
		Data:  items,
		Count: int64(searchResp.Hits.Total.Value),
	}
}

func (r *ProductQueryRepository) buildQuerySearch(dto *productDtos.ProductQueryRequestDto) []map[string]interface{} {
	if dto.Search == "" {
		return nil
	}
	return []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query":  dto.Search,
				"fields": []string{"name^2", "description"},
			},
		},
	}
}

func (r *ProductQueryRepository) buildQueryFilter(dto *productDtos.ProductQueryRequestDto) []map[string]interface{} {
	filters := []map[string]interface{}{}

	if dto.Category != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"category_slug": dto.Category,
			},
		})
	}

	priceRange := map[string]interface{}{}
	if dto.MinPrice > 0 {
		priceRange["gte"] = dto.MinPrice
	}
	if dto.MaxPrice > 0 {
		priceRange["lte"] = dto.MaxPrice
	}
	if len(priceRange) > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}

	return filters
}

func (r *ProductQueryRepository) buildQuery(dto *productDtos.ProductQueryRequestDto) (*bytes.Reader, error) {
	must := r.buildQuerySearch(dto)
	filters := r.buildQueryFilter(dto)

	var query map[string]interface{}
	if len(must) > 0 || len(filters) > 0 {
		boolQuery := map[string]interface{}{}
		if len(must) > 0 {
			boolQuery["must"] = must
		}
		if len(filters) > 0 {
			boolQuery["filter"] = filters
		}
		query = map[string]interface{}{
			"bool": boolQuery,
		}
	} else {
		query = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	body := map[string]interface{}{
		"from": (dto.Page - 1) * dto.PerPage,
		"size": dto.PerPage,
		"sort": []map[string]interface{}{
			{r.resolveSortField(dto.SortBy): map[string]string{"order": string(dto.Order)}},
		},
		"query": query,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	return bytes.NewReader(bodyBytes), nil
}

func (r *ProductQueryRepository) resolveSortField(sortBy string) string {
	switch sortBy {
	case "name":
		return "name.keyword"
	case "price":
		return "price"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
