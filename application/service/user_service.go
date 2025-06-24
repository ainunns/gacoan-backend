package service

import (
	"context"
	"errors"
	"fp-kpl/application/request"
	"fp-kpl/application/response"
	"fp-kpl/domain/user"
	"fp-kpl/infrastructure/database/validation"
	"gorm.io/gorm"
)

type (
	UserService interface {
		Register(ctx context.Context, req request.UserRegister) (response.UserRegister, error)
		GetUserByID(ctx context.Context, userID string) (response.User, error)
		GetUserByEmail(ctx context.Context, email string) (response.User, error)
		Verify(ctx context.Context, req request.UserLogin) (response.RefreshToken, error)
	}

	userService struct {
		userRepository user.Repository
		jwtService     JWTService
		transaction    interface{}
	}
)

func NewUserService(
	userRepository user.Repository,
	jwtService JWTService,
	transaction interface{},
) UserService {
	return &userService{
		userRepository: userRepository,
		jwtService:     jwtService,
		transaction:    transaction,
	}
}

func (s *userService) Register(ctx context.Context, req request.UserRegister) (response.UserRegister, error) {
	_, alreadyExists, err := s.userRepository.CheckEmail(ctx, nil, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response.UserRegister{}, err
	}

	if alreadyExists {
		return response.UserRegister{}, user.ErrorEmailAlreadyExists
	}

	password, err := user.NewPassword(req.Password)
	if err != nil {
		return response.UserRegister{}, err
	}
	role, err := user.NewRole(user.RoleCustomer)
	if err != nil {
		return response.UserRegister{}, err
	}

	userEntity := user.User{
		Email:       req.Email,
		Password:    password,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Role:        role,
	}

	registeredUser, err := s.userRepository.Register(ctx, nil, userEntity)
	if err != nil {
		return response.UserRegister{}, user.ErrorCreateUser
	}

	return response.UserRegister{
		ID:          registeredUser.ID.String(),
		Email:       registeredUser.Email,
		Name:        registeredUser.Name,
		PhoneNumber: registeredUser.PhoneNumber,
		Role:        registeredUser.Role.Name,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (response.User, error) {
	retrievedUser, err := s.userRepository.GetUserByID(ctx, nil, userID)
	if err != nil {
		return response.User{}, user.ErrorGetUserById
	}

	return response.User{
		ID:          retrievedUser.ID.String(),
		Email:       retrievedUser.Email,
		Name:        retrievedUser.Name,
		PhoneNumber: retrievedUser.PhoneNumber,
		Role:        retrievedUser.Role.Name,
	}, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (response.User, error) {
	retrievedUser, err := s.userRepository.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return response.User{}, user.ErrorGetUserByEmail
	}

	return response.User{
		ID:          retrievedUser.ID.String(),
		Email:       retrievedUser.Email,
		Name:        retrievedUser.Name,
		PhoneNumber: retrievedUser.PhoneNumber,
		Role:        retrievedUser.Role.Name,
	}, nil
}

func (s *userService) Verify(ctx context.Context, req request.UserLogin) (response.RefreshToken, error) {
	validatedTransaction, err := validation.ValidateTransaction(s.transaction)
	if err != nil {
		return response.RefreshToken{}, err
	}

	tx, err := validatedTransaction.Begin(ctx)
	if err != nil {
		return response.RefreshToken{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = RecoveredFromPanic(r)
		}
		validatedTransaction.CommitOrRollback(ctx, tx, err)
	}()

	retrievedUser, err := s.userRepository.GetUserByEmail(ctx, tx, req.Email)
	if err != nil {
		return response.RefreshToken{}, user.ErrorEmailNotFound
	}

	checkPassword, err := retrievedUser.Password.IsPasswordMatch([]byte(req.Password))
	if err != nil || !checkPassword {
		return response.RefreshToken{}, err
	}

	accessToken := s.jwtService.GenerateAccessToken(retrievedUser.ID.String(), retrievedUser.Role.Name)

	return response.RefreshToken{
		AccessToken: accessToken,
	}, nil
}
