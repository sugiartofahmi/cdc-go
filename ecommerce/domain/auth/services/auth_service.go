package services

import (
	"context"

	authConstants "go-service/domain/auth/constants"
	authDtos "go-service/domain/auth/dtos"
	authInterfaces "go-service/domain/auth/interfaces"
	roleConstants "go-service/domain/role/constants"
	"go-service/entities"
	"go-service/infrastructure/exceptions"
	"go-service/infrastructure/utils"
)

type AuthService struct {
	authUserQueryRepository authInterfaces.AuthUserQueryRepositoryInterface
	authUserStoreRepository authInterfaces.AuthUserStoreRepositoryInterface
	authRoleQueryRepository authInterfaces.AuthRoleQueryRepositoryInterface
}

func NewAuthService(
	authUserQueryRepository authInterfaces.AuthUserQueryRepositoryInterface,
	authUserStoreRepository authInterfaces.AuthUserStoreRepositoryInterface,
	authRoleQueryRepository authInterfaces.AuthRoleQueryRepositoryInterface,
) authInterfaces.AuthServiceInterface {
	return &AuthService{
		authUserQueryRepository: authUserQueryRepository,
		authUserStoreRepository: authUserStoreRepository,
		authRoleQueryRepository: authRoleQueryRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, dto *authDtos.AuthLoginRequestDto) *authDtos.AuthLoginResponseDto {
	user := s.authUserQueryRepository.FindOneByEmailWithRole(ctx, dto.Email)
	isUserMissing := user == nil
	if isUserMissing {
		panic(*exceptions.UnauthenticatedException(authConstants.AUTH_CREDENTIAL_NOT_VALID))
	}

	isPasswordMismatch := !utils.ComparePassword(user.Password, dto.Password)
	if isPasswordMismatch {
		panic(*exceptions.UnauthenticatedException(authConstants.AUTH_CREDENTIAL_NOT_VALID))
	}

	token, expiresAt := utils.GenerateToken(utils.JWTUser{
		Id:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		RoleId:   user.RoleId,
		RoleName: user.Role.Name,
	})

	return &authDtos.AuthLoginResponseDto{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
	}
}

func (s *AuthService) Register(ctx context.Context, dto *authDtos.AuthRegisterRequestDto) *authDtos.AuthRegisterResponseDto {
	isEmailAlreadyUsed := s.authUserQueryRepository.IsExistsByEmail(ctx, dto.Email)
	if isEmailAlreadyUsed {
		panic(*exceptions.UnprocessableEntityException(authConstants.AUTH_EMAIL_ALREADY_EXISTS))
	}

	developerRole := s.authRoleQueryRepository.FindOneByName(ctx, roleConstants.DEVELOPER)
	isRoleMissing := developerRole == nil
	if isRoleMissing {
		panic(*exceptions.UnprocessableEntityException(authConstants.AUTH_ROLE_NOT_FOUND))
	}

	user := &entities.UserEntity{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		RoleId:   developerRole.Id,
	}
	result := s.authUserStoreRepository.Create(ctx, user)

	return &authDtos.AuthRegisterResponseDto{
		Id:        result.Id,
		Name:      result.Name,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
	}
}
