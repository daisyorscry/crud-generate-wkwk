package repository

import (
	"context"
	"daisy/domain/user/dao"
	"daisy/pkg/database"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*dao.User, error)
	FindAll(ctx context.Context) ([]dao.User, error)
	Paginate(ctx context.Context, limit, offset int) ([]dao.User, int64, error)
	Create(ctx context.Context, data *dao.User) error
	Update(ctx context.Context, data *dao.User) error
	Delete(ctx context.Context, id uint) error
    RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type userRepository struct {
	db database.Connection
}

func New(db database.Connection) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.Run(ctx, fn)
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*dao.User, error) {
	var data dao.User
	err := r.db.Get(ctx).First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]dao.User, error) {
	var list []dao.User
	err := r.db.Get(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *userRepository) Paginate(ctx context.Context, limit, offset int) ([]dao.User, int64, error) {
	var list []dao.User
	var total int64

	tx := r.db.Get(ctx)
	if err := tx.Model(&dao.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := tx.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

func (r *userRepository) Create(ctx context.Context, data *dao.User) error {
	return r.db.Get(ctx).Create(data).Error
}

func (r *userRepository) Update(ctx context.Context, data *dao.User) error {
	return r.db.Get(ctx).Save(data).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Get(ctx).Delete(&dao.User{}, id).Error
}
