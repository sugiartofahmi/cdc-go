package repositories

import (
	"context"
	"log"

	"go-service/entities"
	productInterfaces "go-service/domain/product/interfaces"
	productDtos "go-service/domain/product/dtos"
	"go-service/infrastructure/exceptions"
	infradtos "go-service/infrastructure/dtos"
	"go-service/infrastructure/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductQueryRepository struct {
	db           *gorm.DB
	productModel *gorm.DB
}

func NewProductQueryRepository(db *gorm.DB) productInterfaces.ProductQueryRepositoryInterface {
	return &ProductQueryRepository{
		db:           db,
		productModel: db.Model(&entities.ProductEntity{}),
	}
}

func (r *ProductQueryRepository) Pagination(ctx context.Context, dto *productDtos.ProductQueryRequestDto) *infradtos.PaginationResultDto[entities.ProductEntity] {
	var result infradtos.PaginationResultDto[entities.ProductEntity]

	query := r.productModel.WithContext(ctx).Preload("Category")

	r.querySearch(&query, dto)
	r.querySort(&query, dto)

	err := query.Count(&result.Count).Error
	if err != nil {
		log.Println("Error count products:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	err = query.Scopes(utils.Paginate(&dto.PaginationQueryRequestDto)).Find(&result.Data).Error
	if err != nil {
		log.Println("Error paginate products:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return &result
}

func (r *ProductQueryRepository) FindOneById(ctx context.Context, id uuid.UUID) *entities.ProductEntity {
	var result entities.ProductEntity
	err := r.productModel.WithContext(ctx).Preload("Category").Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		log.Println("Error find product by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &result
}

func (r *ProductQueryRepository) IsExistBySlug(ctx context.Context, slug string) bool {
	var count int64
	err := r.productModel.WithContext(ctx).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		log.Println("Error count product by slug:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *ProductQueryRepository) IsExistBySlugExcludeId(ctx context.Context, slug string, excludeId uuid.UUID) bool {
	var count int64
	err := r.productModel.WithContext(ctx).Where("slug = ? AND id != ?", slug, excludeId).Count(&count).Error
	if err != nil {
		log.Println("Error count product by slug exclude id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *ProductQueryRepository) IsExistById(ctx context.Context, id uuid.UUID) bool {
	var count int64
	err := r.productModel.WithContext(ctx).Where("id = ?", id).Count(&count).Error
	if err != nil {
		log.Println("Error count product by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *ProductQueryRepository) querySearch(query **gorm.DB, dto *productDtos.ProductQueryRequestDto) {
	if dto.Search != "" {
		*query = (*query).Where("name ILIKE ? OR description ILIKE ?", "%"+dto.Search+"%", "%"+dto.Search+"%")
	}
}

func (r *ProductQueryRepository) querySort(query **gorm.DB, dto *productDtos.ProductQueryRequestDto) {
	sortBy := "created_at"
	sortableColumns := []string{"name", "price", "stock", "created_at", "updated_at"}

	if dto.SortBy != "" {
		isColumnAllowed := utils.Contains(sortableColumns, dto.SortBy)
		if isColumnAllowed {
			sortBy = dto.SortBy
		}
	}

	order := "desc"
	if dto.Order != "" {
		order = string(dto.Order)
	}

	*query = (*query).Order(sortBy + " " + order)
}
