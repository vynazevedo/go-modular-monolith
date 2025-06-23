package infra

import (
	"context"
	"errors"
	"time"

	"github.com/vynazevedo/go-modular-monolith/internal/modules/user/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        string `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex"`
	Name      string
	Status    string
	CreatedAt int64
}

func (UserModel) TableName() string {
	return "users"
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(ctx context.Context, user *domain.User) error {
	model := UserModel{
		ID:        user.ID(),
		Email:     user.Email(),
		Name:      user.Name(),
		Status:    user.Status().String(),
		CreatedAt: user.CreatedAt().Unix(),
	}

	result := r.db.WithContext(ctx).Save(&model)
	return result.Error
}

func (r *GormUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var model UserModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	createdAt := time.Unix(model.CreatedAt, 0)

	return domain.ReconstructUser(
		model.ID,
		model.Email,
		model.Name,
		model.Status,
		createdAt,
	)
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel
	result := r.db.WithContext(ctx).First(&model, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	createdAt := time.Unix(model.CreatedAt, 0)

	return domain.ReconstructUser(
		model.ID,
		model.Email,
		model.Name,
		model.Status,
		createdAt,
	)
}

func (r *GormUserRepository) FindAll(ctx context.Context, page, limit int) ([]*domain.User, error) {
	var models []UserModel
	offset := (page - 1) * limit

	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		createdAt := time.Unix(model.CreatedAt, 0)

		user, err := domain.ReconstructUser(
			model.ID,
			model.Email,
			model.Name,
			model.Status,
			createdAt,
		)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return users, nil
}

func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
