package repository

import (
	"context"
	"daisy/domain/todo/dao"
	"daisy/pkg/database"
)

type TodoRepository interface {
	FindByID(ctx context.Context, id uint) (*dao.Todo, error)
	FindAll(ctx context.Context) ([]dao.Todo, error)
	Paginate(ctx context.Context, limit, offset int) ([]dao.Todo, int64, error)
	Create(ctx context.Context, data *dao.Todo) error
	Update(ctx context.Context, data *dao.Todo) error
	Delete(ctx context.Context, id uint) error
    RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type todoRepository struct {
	db database.Connection
}

func New(db database.Connection) TodoRepository {
	return &todoRepository{
		db: db,
	}
}

func (r *todoRepository) RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.Run(ctx, fn)
}

func (r *todoRepository) FindByID(ctx context.Context, id uint) (*dao.Todo, error) {
	var data dao.Todo
	err := r.db.Get(ctx).First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *todoRepository) FindAll(ctx context.Context) ([]dao.Todo, error) {
	var list []dao.Todo
	err := r.db.Get(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *todoRepository) Paginate(ctx context.Context, limit, offset int) ([]dao.Todo, int64, error) {
	var list []dao.Todo
	var total int64

	tx := r.db.Get(ctx)
	if err := tx.Model(&dao.Todo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := tx.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

func (r *todoRepository) Create(ctx context.Context, data *dao.Todo) error {
	return r.db.Get(ctx).Create(data).Error
}

func (r *todoRepository) Update(ctx context.Context, data *dao.Todo) error {
	return r.db.Get(ctx).Save(data).Error
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Get(ctx).Delete(&dao.Todo{}, id).Error
}
