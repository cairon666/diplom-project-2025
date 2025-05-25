package user_service

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

type UserRepo interface {
	GetUserById(ctx context.Context, ID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
}

type UserService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (userService *UserService) GetUserById(ctx context.Context, ID uuid.UUID) (*models.User, error) {
	return userService.userRepo.GetUserById(ctx, ID)
}

func (userService *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return userService.userRepo.GetUserByEmail(ctx, email)
}

func (userService *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return userService.userRepo.CreateUser(ctx, user)
}

func (userService *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return userService.userRepo.UpdateUser(ctx, user)
}
