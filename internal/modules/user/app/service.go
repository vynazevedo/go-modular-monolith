package app

import (
	"context"
	"errors"

	"github.com/vynazevedo/template-go-modular/internal/modules/user/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, cmd CreateUserCommand) (*domain.UserInfo, error) {
	existingUser, err := s.repo.FindByEmail(ctx, cmd.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	user, err := domain.NewUser(cmd.Email, cmd.Name)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:     user.ID(),
		Email:  user.Email(),
		Name:   user.Name(),
		Status: user.Status().String(),
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, query GetUserQuery) (*domain.UserInfo, error) {
	user, err := s.repo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:     user.ID(),
		Email:  user.Email(),
		Name:   user.Name(),
		Status: user.Status().String(),
	}, nil
}

func (s *UserService) QueryUserByEmail(ctx context.Context, query GetUserByEmailQuery) (*domain.UserInfo, error) {
	user, err := s.repo.FindByEmail(ctx, query.Email)
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:     user.ID(),
		Email:  user.Email(),
		Name:   user.Name(),
		Status: user.Status().String(),
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, cmd UpdateUserCommand) (*domain.UserInfo, error) {
	user, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := user.UpdateName(cmd.Name); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:     user.ID(),
		Email:  user.Email(),
		Name:   user.Name(),
		Status: user.Status().String(),
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, cmd DeleteUserCommand) error {
	return s.repo.Delete(ctx, cmd.ID)
}

func (s *UserService) ActivateUser(ctx context.Context, cmd ActivateUserCommand) (*domain.UserInfo, error) {
	user, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	user.Activate()

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:     user.ID(),
		Email:  user.Email(),
		Name:   user.Name(),
		Status: user.Status().String(),
	}, nil
}

func (s *UserService) DeactivateUser(ctx context.Context, cmd DeactivateUserCommand) (*domain.UserInfo, error) {
	user, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	user.Deactivate()

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:     user.ID(),
		Email:  user.Email(),
		Name:   user.Name(),
		Status: user.Status().String(),
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, query ListUsersQuery) ([]*domain.UserInfo, error) {
	users, err := s.repo.FindAll(ctx, query.Page, query.Limit)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.UserInfo, len(users))
	for i, user := range users {
		result[i] = &domain.UserInfo{
			ID:     user.ID(),
			Email:  user.Email(),
			Name:   user.Name(),
			Status: user.Status().String(),
		}
	}

	return result, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, id string) (*domain.UserInfo, error) {
	return s.GetUser(ctx, GetUserQuery{ID: id})
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.UserInfo, error) {
	return s.QueryUserByEmail(ctx, GetUserByEmailQuery{Email: email})
}
