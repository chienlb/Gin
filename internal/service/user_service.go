package service

import (
	"gin-demo/internal/domain"
	"gin-demo/internal/repository"
	"gin-demo/pkg/apperror"
	"gin-demo/pkg/logger"
	"gin-demo/pkg/utils"
	"gin-demo/pkg/validator"
)

type UserService struct {
	repo      *repository.UserRepository
	validator *validator.UserValidator
	log       *logger.Logger
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo:      repo,
		validator: validator.NewUserValidator(),
		log:       logger.Get(),
	}
}

func (s *UserService) CreateUser(req *domain.CreateUserRequest) (*domain.UserResponse, *apperror.AppError) {
	// Validate input
	if err := s.validator.ValidateCreateRequest(req.Name, req.Email, req.Password); err != nil {
		return nil, err
	}

	// Normalize email
	normalizedEmail := utils.NormalizeEmail(req.Email)

	// Check if user already exists
	existingUser, _ := s.repo.GetByEmail(normalizedEmail)
	if existingUser != nil {
		return nil, apperror.DuplicateEmailError(normalizedEmail)
	}

	// Hash password
	hashedPassword := utils.HashPassword(req.Password)

	// Create user
	user := &domain.User{
		Name:     utils.TrimSpaces(req.Name),
		Email:    normalizedEmail,
		Password: hashedPassword,
	}

	if err := s.repo.Create(user); err != nil {
		s.log.Error("Failed to create user", err)
		return nil, apperror.NewWithError(
			apperror.CodeInternalServerError,
			"Failed to create user",
			500,
			err,
		)
	}

	s.log.Info("User created successfully: " + normalizedEmail)
	return userToResponse(user), nil
}

func (s *UserService) GetUser(id int) (*domain.UserResponse, *apperror.AppError) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.UserNotFoundError(id)
	}

	return userToResponse(user), nil
}

func (s *UserService) GetAllUsers() ([]domain.UserResponse, *apperror.AppError) {
	users, err := s.repo.GetAll()
	if err != nil {
		s.log.Error("Failed to get users", err)
		return nil, apperror.NewWithError(
			apperror.CodeInternalServerError,
			"Failed to retrieve users",
			500,
			err,
		)
	}

	var responses []domain.UserResponse
	for _, user := range users {
		responses = append(responses, *userToResponse(&user))
	}

	return responses, nil
}

func (s *UserService) UpdateUser(id int, req *domain.UpdateUserRequest) (*domain.UserResponse, *apperror.AppError) {
	// Validate input
	if err := s.validator.ValidateUpdateRequest(req.Name, req.Email); err != nil {
		return nil, err
	}

	// Get existing user
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.UserNotFoundError(id)
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = utils.TrimSpaces(req.Name)
	}
	if req.Email != "" {
		normalizedEmail := utils.NormalizeEmail(req.Email)
		// Check if email is already taken by another user
		existingUser, _ := s.repo.GetByEmail(normalizedEmail)
		if existingUser != nil && existingUser.ID != id {
			return nil, apperror.DuplicateEmailError(normalizedEmail)
		}
		user.Email = normalizedEmail
	}

	if err := s.repo.Update(user); err != nil {
		s.log.Error("Failed to update user", err)
		return nil, apperror.NewWithError(
			apperror.CodeInternalServerError,
			"Failed to update user",
			500,
			err,
		)
	}

	s.log.Info("User updated successfully: ID " + string(rune(user.ID)))
	return userToResponse(user), nil
}

func (s *UserService) DeleteUser(id int) *apperror.AppError {
	err := s.repo.Delete(id)
	if err != nil {
		return apperror.UserNotFoundError(id)
	}

	s.log.Info("User deleted successfully: ID " + string(rune(id)))
	return nil
}

func userToResponse(user *domain.User) *domain.UserResponse {
	return &domain.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
