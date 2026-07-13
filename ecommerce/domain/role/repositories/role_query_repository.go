package repositories

import (
	"context"
	"log"

	"go-service/entities"
	roleInterfaces "go-service/domain/role/interfaces"
	roleDtos "go-service/domain/role/dtos"
	"go-service/infrastructure/exceptions"
	infradtos "go-service/infrastructure/dtos"
	"go-service/infrastructure/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleQueryRepository struct {
	db        *gorm.DB
	roleModel *gorm.DB
}

func NewRoleQueryRepository(db *gorm.DB) roleInterfaces.RoleQueryRepositoryInterface {
	return &RoleQueryRepository{
		db:        db,
		roleModel: db.Model(&entities.RoleEntity{}),
	}
}

func (r *RoleQueryRepository) Pagination(ctx context.Context, dto *roleDtos.RoleQueryRequestDto) *infradtos.PaginationResultDto[entities.RoleEntity] {
	var result infradtos.PaginationResultDto[entities.RoleEntity]

	query := r.roleModel.WithContext(ctx)

	r.querySearch(&query, dto)
	r.querySort(&query, dto)

	err := query.Count(&result.Count).Error
	if err != nil {
		log.Println("Error count roles:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	err = query.Scopes(utils.Paginate(&dto.PaginationQueryRequestDto)).Find(&result.Data).Error
	if err != nil {
		log.Println("Error paginate roles:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return &result
}

func (r *RoleQueryRepository) FindOneById(ctx context.Context, id uuid.UUID) *entities.RoleEntity {
	var result entities.RoleEntity
	err := r.roleModel.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		log.Println("Error find role by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &result
}

func (r *RoleQueryRepository) FindOneByName(ctx context.Context, name string) *entities.RoleEntity {
	var result entities.RoleEntity
	err := r.roleModel.WithContext(ctx).Where("name = ?", name).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		log.Println("Error find role by name:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &result
}

func (r *RoleQueryRepository) IsExistByName(ctx context.Context, name string) bool {
	var count int64
	err := r.roleModel.WithContext(ctx).Where("name = ?", name).Count(&count).Error
	if err != nil {
		log.Println("Error count role by name:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *RoleQueryRepository) IsExistByNameExcludeId(ctx context.Context, name string, excludeId uuid.UUID) bool {
	var count int64
	err := r.roleModel.WithContext(ctx).Where("name = ? AND id != ?", name, excludeId).Count(&count).Error
	if err != nil {
		log.Println("Error count role by name exclude id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *RoleQueryRepository) IsExistById(ctx context.Context, id uuid.UUID) bool {
	var count int64
	err := r.roleModel.WithContext(ctx).Where("id = ?", id).Count(&count).Error
	if err != nil {
		log.Println("Error count role by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *RoleQueryRepository) querySearch(query **gorm.DB, dto *roleDtos.RoleQueryRequestDto) {
	if dto.Search != "" {
		*query = (*query).Where("name ILIKE ?", "%"+dto.Search+"%")
	}
}

func (r *RoleQueryRepository) querySort(query **gorm.DB, dto *roleDtos.RoleQueryRequestDto) {
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
