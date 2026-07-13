package repositories

import (
	"context"
	"log"

	"go-service/entities"
	categoryInterfaces "go-service/domain/category/interfaces"
	categoryDtos "go-service/domain/category/dtos"
	"go-service/infrastructure/exceptions"
	infradtos "go-service/infrastructure/dtos"
	"go-service/infrastructure/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryQueryRepository struct {
	db            *gorm.DB
	categoryModel *gorm.DB
}

func NewCategoryQueryRepository(db *gorm.DB) categoryInterfaces.CategoryQueryRepositoryInterface {
	return &CategoryQueryRepository{
		db:            db,
		categoryModel: db.Model(&entities.CategoryEntity{}),
	}
}

func (r *CategoryQueryRepository) Pagination(ctx context.Context, dto *categoryDtos.CategoryQueryRequestDto) *infradtos.PaginationResultDto[entities.CategoryEntity] {
	var result infradtos.PaginationResultDto[entities.CategoryEntity]

	query := r.categoryModel.WithContext(ctx)

	r.querySearch(&query, dto)
	r.querySort(&query, dto)

	err := query.Count(&result.Count).Error
	if err != nil {
		log.Println("Error count categories:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	err = query.Scopes(utils.Paginate(&dto.PaginationQueryRequestDto)).Find(&result.Data).Error
	if err != nil {
		log.Println("Error paginate categories:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return &result
}

func (r *CategoryQueryRepository) FindOneById(ctx context.Context, id uuid.UUID) *entities.CategoryEntity {
	var result entities.CategoryEntity
	err := r.categoryModel.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		log.Println("Error find category by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &result
}

func (r *CategoryQueryRepository) FindOneBySlug(ctx context.Context, slug string) *entities.CategoryEntity {
	var result entities.CategoryEntity
	err := r.categoryModel.WithContext(ctx).Where("slug = ?", slug).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		log.Println("Error find category by slug:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &result
}

func (r *CategoryQueryRepository) IsExistByName(ctx context.Context, name string) bool {
	var count int64
	err := r.categoryModel.WithContext(ctx).Where("name = ?", name).Count(&count).Error
	if err != nil {
		log.Println("Error count category by name:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *CategoryQueryRepository) IsExistByNameExcludeId(ctx context.Context, name string, excludeId uuid.UUID) bool {
	var count int64
	err := r.categoryModel.WithContext(ctx).Where("name = ? AND id != ?", name, excludeId).Count(&count).Error
	if err != nil {
		log.Println("Error count category by name exclude id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *CategoryQueryRepository) IsExistById(ctx context.Context, id uuid.UUID) bool {
	var count int64
	err := r.categoryModel.WithContext(ctx).Where("id = ?", id).Count(&count).Error
	if err != nil {
		log.Println("Error count category by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *CategoryQueryRepository) querySearch(query **gorm.DB, dto *categoryDtos.CategoryQueryRequestDto) {
	if dto.Search != "" {
		*query = (*query).Where("name ILIKE ?", "%"+dto.Search+"%")
	}
}

func (r *CategoryQueryRepository) querySort(query **gorm.DB, dto *categoryDtos.CategoryQueryRequestDto) {
	sortBy := "created_at"
	sortableColumns := []string{"name", "created_at", "updated_at"}

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
