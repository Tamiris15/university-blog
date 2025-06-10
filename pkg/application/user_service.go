package application

import (
	"context"
	"errors"

	"github.com/Tamiris15/university-blog/pkg/domain"
)

type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("Недопустимые или отсутствующие данные пользователя")
	}

	if !IsValidEmail(user.Email) {
		return ErrInvalidEmail
	}

	existingUser, err := s.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrEmailAlreadyUsed
	}

	return s.userRepository.Create(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepository.GetByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.ID == 0 {
		return errors.New("Недопустимый идентификатор пользователя")
	}
	return s.userRepository.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepository.Delete(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepository.GetAll(ctx)
}
