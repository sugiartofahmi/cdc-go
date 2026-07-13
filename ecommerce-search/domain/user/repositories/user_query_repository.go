package repositories

import (
	"context"
	"log"

	"go-service/entities"
	userInterfaces "go-service/domain/user/interfaces"
	userDtos "go-service/domain/user/dtos"
	"go-service/infrastructure/exceptions"
	infradtos "go-service/infrastructure/dtos"
	"go-service/infrastructure/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserQueryRepository struct {
	db        *gorm.DB
	userModel *gorm.DB
}

func NewUserQueryRepository(db *gorm.DB) userInterfaces.UserQueryRepositoryInterface {
	return &UserQueryRepository{
		db:        db,
		userModel: db.Model(&entities.UserEntity{}),
	}
}

func (r *UserQueryRepository) Pagination(ctx context.Context, dto *userDtos.UserQueryRequestDto) *infradtos.PaginationResultDto[entities.UserEntity] {
	var result infradtos.PaginationResultDto[entities.UserEntity]

	query := r.userModel.WithContext(ctx)

	r.querySearch(&query, dto)
	r.querySort(&query, dto)

	err := query.Count(&result.Count).Error
	if err != nil {
		log.Println("Error count users:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	selectColumns := []string{"id", "name", "email", "status", "role_id", "created_at", "updated_at"}
	err = query.Scopes(utils.Paginate(&dto.PaginationQueryRequestDto)).Preload("Role").Select(selectColumns).Find(&result.Data).Error
	if err != nil {
		log.Println("Error paginate users:", err)
		panic(*exceptions.ServerErrorException(err))
	}

	return &result
}

func (r *UserQueryRepository) FindOneById(ctx context.Context, id uuid.UUID) *entities.UserEntity {
	query := r.userModel.WithContext(ctx).Preload("Role")
	var result entities.UserEntity
	err := query.Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		log.Println("Error find user by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &result
}

func (r *UserQueryRepository) FindOneByIdForDetail(ctx context.Context, id uuid.UUID) *entities.UserEntity {
	var result entities.UserEntity
	err := r.userModel.WithContext(ctx).
		Select("id", "name", "email", "status", "role_id", "created_at", "updated_at").
		Preload("Role").
		Where("id = ?", id).First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		log.Println("Error find user by id for detail:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return &result
}

func (r *UserQueryRepository) IsExistByName(ctx context.Context, name string) bool {
	var count int64
	err := r.userModel.WithContext(ctx).Where("name = ?", name).Count(&count).Error
	if err != nil {
		log.Println("Error count user by name:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *UserQueryRepository) IsExistByNameExcludeId(ctx context.Context, name string, excludeId uuid.UUID) bool {
	var count int64
	err := r.userModel.WithContext(ctx).Where("name = ? AND id != ?", name, excludeId).Count(&count).Error
	if err != nil {
		log.Println("Error count user by name exclude id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *UserQueryRepository) IsExistByEmail(ctx context.Context, email string) bool {
	var count int64
	err := r.userModel.WithContext(ctx).Where("email = ?", email).Count(&count).Error
	if err != nil {
		log.Println("Error count user by email:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *UserQueryRepository) IsExistById(ctx context.Context, id uuid.UUID) bool {
	var count int64
	err := r.userModel.WithContext(ctx).Where("id = ?", id).Count(&count).Error
	if err != nil {
		log.Println("Error count user by id:", err)
		panic(*exceptions.ServerErrorException(err))
	}
	return count > 0
}

func (r *UserQueryRepository) querySearch(query **gorm.DB, dto *userDtos.UserQueryRequestDto) {
	if dto.Search != "" {
		*query = (*query).Where("name ILIKE ? OR email ILIKE ?", "%"+dto.Search+"%", "%"+dto.Search+"%")
	}
}

func (r *UserQueryRepository) querySort(query **gorm.DB, dto *userDtos.UserQueryRequestDto) {
	sortBy := "created_at"
	sortableColumns := []string{"name", "email", "created_at", "updated_at"}

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
