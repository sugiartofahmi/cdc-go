package services

import (
	"context"

	"go-service/entities"
	userInterfaces "go-service/domain/user/interfaces"
	userConstants "go-service/domain/user/constants"
	userDtos "go-service/domain/user/dtos"
	"go-service/infrastructure/exceptions"
	infradtos "go-service/infrastructure/dtos"
	"go-service/infrastructure/utils"

	"github.com/google/uuid"
)

type UserService struct {
	userQueryRepository     userInterfaces.UserQueryRepositoryInterface
	userStoreRepository     userInterfaces.UserStoreRepositoryInterface
	userRoleQueryRepository userInterfaces.UserRoleQueryRepositoryInterface
}

func NewUserService(
	userQueryRepository userInterfaces.UserQueryRepositoryInterface,
	userStoreRepository userInterfaces.UserStoreRepositoryInterface,
	userRoleQueryRepository userInterfaces.UserRoleQueryRepositoryInterface,
) userInterfaces.UserServiceInterface {
	return &UserService{
		userQueryRepository:     userQueryRepository,
		userStoreRepository:     userStoreRepository,
		userRoleQueryRepository: userRoleQueryRepository,
	}
}

func (s *UserService) Pagination(ctx context.Context, dto *userDtos.UserQueryRequestDto) *infradtos.PaginationResultDto[entities.UserEntity] {
	return s.userQueryRepository.Pagination(ctx, dto)
}

func (s *UserService) Detail(ctx context.Context, id uuid.UUID) *entities.UserEntity {
	user := s.userQueryRepository.FindOneById(ctx, id)
	isUserNotFound := user == nil
	if isUserNotFound {
		panic(*exceptions.NotFoundException(userConstants.USER_NOT_FOUND))
	}
	return user
}

func (s *UserService) Create(ctx context.Context, dto *userDtos.UserCreateRequestDto) *entities.UserEntity {
	isNameAlreadyUsed := s.userQueryRepository.IsExistByName(ctx, dto.Name)
	if isNameAlreadyUsed {
		panic(*exceptions.ConflictException(userConstants.USER_NAME_ALREADY_EXISTS))
	}
	isEmailAlreadyUsed := s.userQueryRepository.IsExistByEmail(ctx, dto.Email)
	if isEmailAlreadyUsed {
		panic(*exceptions.ConflictException(userConstants.USER_EMAIL_ALREADY_EXISTS))
	}

	isRoleNotFound := !s.userRoleQueryRepository.IsExistById(ctx, dto.RoleId)
	if isRoleNotFound {
		panic(*exceptions.UnprocessableEntityException(userConstants.USER_ROLE_NOT_FOUND))
	}

	user := &entities.UserEntity{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		RoleId:   dto.RoleId,
	}
	return s.userStoreRepository.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, dto *userDtos.UserUpdateRequestDto) *entities.UserEntity {
	user := s.userQueryRepository.FindOneById(ctx, id)
	isUserNotFound := user == nil
	if isUserNotFound {
		panic(*exceptions.NotFoundException(userConstants.USER_NOT_FOUND))
	}

	hasNameUpdate := dto.Name != nil
	if hasNameUpdate {
		isNameAlreadyUsedByOtherUser := s.userQueryRepository.IsExistByNameExcludeId(ctx, *dto.Name, id)
		if isNameAlreadyUsedByOtherUser {
			panic(*exceptions.ConflictException(userConstants.USER_NAME_ALREADY_EXISTS))
		}
		user.Name = *dto.Name
	}
	hasPasswordUpdate := dto.Password != nil
	if hasPasswordUpdate {
		user.Password = utils.HashPassword(*dto.Password)
	}
	hasStatusUpdate := dto.Status != nil
	if hasStatusUpdate {
		user.Status = *dto.Status
	}
	hasRoleIdUpdate := dto.RoleId != nil
	if hasRoleIdUpdate {
		isRoleNotFound := !s.userRoleQueryRepository.IsExistById(ctx, *dto.RoleId)
		if isRoleNotFound {
			panic(*exceptions.UnprocessableEntityException(userConstants.USER_ROLE_NOT_FOUND))
		}
		user.RoleId = *dto.RoleId
	}

	return s.userStoreRepository.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) {
	user := s.userQueryRepository.FindOneById(ctx, id)
	isUserNotFound := user == nil
	if isUserNotFound {
		panic(*exceptions.NotFoundException(userConstants.USER_NOT_FOUND))
	}
	s.userStoreRepository.Delete(ctx, id)
}
