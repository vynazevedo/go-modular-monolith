package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/vynazevedo/go-modular-monolith/internal/modules/user/domain"
)

type MockUserRepository struct {
	users map[string]*domain.User

	SaveFunc        func(ctx context.Context, user *domain.User) error
	FindByIDFunc    func(ctx context.Context, id string) (*domain.User, error)
	FindByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	FindAllFunc     func(ctx context.Context, page, limit int) ([]*domain.User, error)
	DeleteFunc      func(ctx context.Context, id string) error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, user)
	}
	m.users[user.ID()] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	for _, user := range m.users {
		if user.Email() == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindAll(ctx context.Context, page, limit int) ([]*domain.User, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, page, limit)
	}

	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(users) {
		return []*domain.User{}, nil
	}

	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) AddUser(user *domain.User) {
	m.users[user.ID()] = user
}

var _ domain.UserRepository = (*MockUserRepository)(nil)

func TestCreateUser(t *testing.T) {
	t.Run("Create valid user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		cmd := CreateUserCommand{
			Email: "test@example.com",
			Name:  "Test User",
		}

		user, err := service.CreateUser(context.Background(), cmd)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.Email != cmd.Email {
			t.Errorf("Expected email %s, got %s", cmd.Email, user.Email)
		}

		if user.Name != cmd.Name {
			t.Errorf("Expected name %s, got %s", cmd.Name, user.Name)
		}

		if user.Status != domain.StatusActive.String() {
			t.Errorf("Expected status %s, got %s", domain.StatusActive, user.Status)
		}
	})

	t.Run("Create user with existing email", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		existingUser, _ := domain.NewUser("test@example.com", "Existing User")
		repo.FindByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
			if email == "test@example.com" {
				return existingUser, nil
			}
			return nil, errors.New("user not found")
		}

		cmd := CreateUserCommand{
			Email: "test@example.com",
			Name:  "Test User",
		}

		user, err := service.CreateUser(context.Background(), cmd)

		if err == nil {
			t.Error("Expected error for duplicate email, got nil")
		}

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})

	t.Run("Create user with invalid data", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		cmd := CreateUserCommand{
			Email: "",
			Name:  "",
		}

		user, err := service.CreateUser(context.Background(), cmd)

		if err == nil {
			t.Error("Expected error for invalid data, got nil")
		}

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})

	t.Run("Create user with repository error", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		repo.SaveFunc = func(ctx context.Context, user *domain.User) error {
			return errors.New("database error")
		}

		cmd := CreateUserCommand{
			Email: "test@example.com",
			Name:  "Test User",
		}

		user, err := service.CreateUser(context.Background(), cmd)

		if err == nil {
			t.Error("Expected error from repository, got nil")
		}

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Get existing user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		testUser, _ := domain.NewUser("test@example.com", "Test User")
		repo.AddUser(testUser)

		query := GetUserQuery{ID: testUser.ID()}

		user, err := service.GetUser(context.Background(), query)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.ID != testUser.ID() {
			t.Errorf("Expected ID %s, got %s", testUser.ID(), user.ID)
		}

		if user.Email != testUser.Email() {
			t.Errorf("Expected email %s, got %s", testUser.Email(), user.Email)
		}

		if user.Name != testUser.Name() {
			t.Errorf("Expected name %s, got %s", testUser.Name(), user.Name)
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		query := GetUserQuery{ID: "non-existent-id"}

		user, err := service.GetUser(context.Background(), query)

		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Update existing user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		testUser, _ := domain.NewUser("test@example.com", "Original Name")
		repo.AddUser(testUser)

		cmd := UpdateUserCommand{
			ID:   testUser.ID(),
			Name: "Updated Name",
		}

		user, err := service.UpdateUser(context.Background(), cmd)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.Name != cmd.Name {
			t.Errorf("Expected name %s, got %s", cmd.Name, user.Name)
		}

		updatedUser, _ := repo.FindByID(context.Background(), testUser.ID())
		if updatedUser.Name() != cmd.Name {
			t.Errorf("User not updated in repository. Expected name %s, got %s", cmd.Name, updatedUser.Name())
		}
	})

	t.Run("Update non-existent user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		cmd := UpdateUserCommand{
			ID:   "non-existent-id",
			Name: "Updated Name",
		}

		user, err := service.UpdateUser(context.Background(), cmd)

		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})

	t.Run("Update with invalid name", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		testUser, _ := domain.NewUser("test@example.com", "Original Name")
		repo.AddUser(testUser)

		cmd := UpdateUserCommand{
			ID:   testUser.ID(),
			Name: "",
		}

		user, err := service.UpdateUser(context.Background(), cmd)

		if err == nil {
			t.Error("Expected error for invalid name, got nil")
		}

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Delete existing user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		testUser, _ := domain.NewUser("test@example.com", "Test User")
		repo.AddUser(testUser)

		cmd := DeleteUserCommand{ID: testUser.ID()}

		err := service.DeleteUser(context.Background(), cmd)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		_, err = repo.FindByID(context.Background(), testUser.ID())
		if err == nil {
			t.Error("Expected user to be deleted, but it still exists")
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		cmd := DeleteUserCommand{ID: "non-existent-id"}

		err := service.DeleteUser(context.Background(), cmd)

		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
	})
}

func TestActivateDeactivateUser(t *testing.T) {
	t.Run("Activate user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		testUser, _ := domain.NewUser("test@example.com", "Test User")
		testUser.Deactivate()
		repo.AddUser(testUser)

		cmd := ActivateUserCommand{ID: testUser.ID()}

		user, err := service.ActivateUser(context.Background(), cmd)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.Status != domain.StatusActive.String() {
			t.Errorf("Expected status %s, got %s", domain.StatusActive, user.Status)
		}

		updatedUser, _ := repo.FindByID(context.Background(), testUser.ID())
		if updatedUser.Status() != domain.StatusActive {
			t.Errorf("User not activated in repository. Expected status %s, got %s", domain.StatusActive, updatedUser.Status())
		}
	})

	t.Run("Deactivate user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		testUser, _ := domain.NewUser("test@example.com", "Test User")
		repo.AddUser(testUser)

		cmd := DeactivateUserCommand{ID: testUser.ID()}

		user, err := service.DeactivateUser(context.Background(), cmd)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.Status != domain.StatusInactive.String() {
			t.Errorf("Expected status %s, got %s", domain.StatusInactive, user.Status)
		}

		updatedUser, _ := repo.FindByID(context.Background(), testUser.ID())
		if updatedUser.Status() != domain.StatusInactive {
			t.Errorf("User not deactivated in repository. Expected status %s, got %s", domain.StatusInactive, updatedUser.Status())
		}
	})

	t.Run("Activate non-existent user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		cmd := ActivateUserCommand{ID: "non-existent-id"}

		user, err := service.ActivateUser(context.Background(), cmd)

		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})

	t.Run("Deactivate non-existent user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		cmd := DeactivateUserCommand{ID: "non-existent-id"}

		user, err := service.DeactivateUser(context.Background(), cmd)

		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}

		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})
}

func TestListUsers(t *testing.T) {
	t.Run("List users with pagination", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		for i := 0; i < 15; i++ {
			user, _ := domain.NewUser(
				fmt.Sprintf("user%d@example.com", i),
				fmt.Sprintf("User %d", i),
			)
			repo.AddUser(user)
		}

		query1 := ListUsersQuery{
			Page:  1,
			Limit: 5,
		}

		users1, err := service.ListUsers(context.Background(), query1)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(users1) != 5 {
			t.Errorf("Expected 5 users on first page, got %d", len(users1))
		}

		query2 := ListUsersQuery{
			Page:  2,
			Limit: 5,
		}

		users2, err := service.ListUsers(context.Background(), query2)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(users2) != 5 {
			t.Errorf("Expected 5 users on second page, got %d", len(users2))
		}

		query3 := ListUsersQuery{
			Page:  3,
			Limit: 5,
		}

		users3, err := service.ListUsers(context.Background(), query3)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(users3) != 5 {
			t.Errorf("Expected 5 users on third page, got %d", len(users3))
		}

		query4 := ListUsersQuery{
			Page:  4,
			Limit: 5,
		}

		users4, err := service.ListUsers(context.Background(), query4)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(users4) != 0 {
			t.Errorf("Expected 0 users on fourth page, got %d", len(users4))
		}
	})

	t.Run("List users with empty repository", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		query := ListUsersQuery{
			Page:  1,
			Limit: 10,
		}

		users, err := service.ListUsers(context.Background(), query)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(users) != 0 {
			t.Errorf("Expected empty list, got %d users", len(users))
		}
	})

	t.Run("List users with repository error", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		repo.FindAllFunc = func(ctx context.Context, page, limit int) ([]*domain.User, error) {
			return nil, errors.New("database error")
		}

		query := ListUsersQuery{
			Page:  1,
			Limit: 10,
		}

		users, err := service.ListUsers(context.Background(), query)

		if err == nil {
			t.Error("Expected error from repository, got nil")
		}

		if users != nil {
			t.Errorf("Expected nil users, got %v", users)
		}
	})
}
