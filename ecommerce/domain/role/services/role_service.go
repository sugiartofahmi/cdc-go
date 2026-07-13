package services

import (
	"context"

	"go-service/entities"
	roleInterfaces "go-service/domain/role/interfaces"
	roleConstants "go-service/domain/role/constants"
	roleDtos "go-service/domain/role/dtos"
	"go-service/infrastructure/exceptions"
	infradtos "go-service/infrastructure/dtos"

	"github.com/google/uuid"
)

type RoleService struct {
	roleQueryRepository roleInterfaces.RoleQueryRepositoryInterface
	roleStoreRepository roleInterfaces.RoleStoreRepositoryInterface
}

func NewRoleService(
	roleQueryRepository roleInterfaces.RoleQueryRepositoryInterface,
	roleStoreRepository roleInterfaces.RoleStoreRepositoryInterface,
) roleInterfaces.RoleServiceInterface {
	return &RoleService{
		roleQueryRepository: roleQueryRepository,
		roleStoreRepository: roleStoreRepository,
	}
}

func (s *RoleService) Pagination(ctx context.Context, dto *roleDtos.RoleQueryRequestDto) *infradtos.PaginationResultDto[entities.RoleEntity] {
	return s.roleQueryRepository.Pagination(ctx, dto)
}

func (s *RoleService) Detail(ctx context.Context, id uuid.UUID) *entities.RoleEntity {
	role := s.roleQueryRepository.FindOneById(ctx, id)
	isRoleNotFound := role == nil
	if isRoleNotFound {
		panic(*exceptions.NotFoundException(roleConstants.ROLE_NOT_FOUND))
	}
	return role
}

func (s *RoleService) Create(ctx context.Context, dto *roleDtos.RoleCreateRequestDto) *entities.RoleEntity {
	isNameAlreadyUsed := s.roleQueryRepository.IsExistByName(ctx, dto.Name)
	if isNameAlreadyUsed {
		panic(*exceptions.ConflictException(roleConstants.ROLE_NAME_ALREADY_EXISTS))
	}

	entity := &entities.RoleEntity{
		Name: dto.Name,
	}
	return s.roleStoreRepository.Create(ctx, entity)
}

func (s *RoleService) Update(ctx context.Context, id uuid.UUID, dto *roleDtos.RoleUpdateRequestDto) *entities.RoleEntity {
	role := s.roleQueryRepository.FindOneById(ctx, id)
	isRoleNotFound := role == nil
	if isRoleNotFound {
		panic(*exceptions.NotFoundException(roleConstants.ROLE_NOT_FOUND))
	}

	hasNameUpdate := dto.Name != nil
	if hasNameUpdate {
		isNameAlreadyUsedByOtherRole := s.roleQueryRepository.IsExistByNameExcludeId(ctx, *dto.Name, id)
		if isNameAlreadyUsedByOtherRole {
			panic(*exceptions.ConflictException(roleConstants.ROLE_NAME_ALREADY_EXISTS))
		}
		role.Name = *dto.Name
	}

	return s.roleStoreRepository.Update(ctx, role)
}

func (s *RoleService) Delete(ctx context.Context, id uuid.UUID) {
	role := s.roleQueryRepository.FindOneById(ctx, id)
	isRoleNotFound := role == nil
	if isRoleNotFound {
		panic(*exceptions.NotFoundException(roleConstants.ROLE_NOT_FOUND))
	}
	s.roleStoreRepository.Delete(ctx, id)
}
