package repository

import (
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"WalletTopUp/internal/domain/repository"
	"WalletTopUp/internal/infra/database/mapper"
	"WalletTopUp/internal/infra/database/model"
	"WalletTopUp/pkg/lib"
	"context"
	"errors"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(logger lib.Logger, db *gorm.DB) repository.UserRepo {
	return &userRepo{
		db,
	}
}

func (repo *userRepo) Create(user entity.User) error {
	model := mapper.ToUserModel(user)
	return repo.db.Create(&model).Error
}

func (repo *userRepo) FindById(ctx context.Context, id uint) (*entity.User, error) {
	var user model.User
	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.ErrUserNotFound
		}
		return nil, err
	}
	return user.ToEntity(), nil
}
