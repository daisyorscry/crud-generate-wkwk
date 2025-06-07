package repository

import (
	"context"
	"daisy/domain/article/dao"
	"daisy/pkg/database"
)

type ArticleRepository interface {
	FindByID(ctx context.Context, id uint) (*dao.Article, error)
	FindAll(ctx context.Context) ([]dao.Article, error)
	Paginate(ctx context.Context, limit, offset int) ([]dao.Article, int64, error)
	Create(ctx context.Context, data *dao.Article) error
	Update(ctx context.Context, data *dao.Article) error
	Delete(ctx context.Context, id uint) error
    RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type articleRepository struct {
	db database.Connection
}

func New(db database.Connection) ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) RunTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.Run(ctx, fn)
}

func (r *articleRepository) FindByID(ctx context.Context, id uint) (*dao.Article, error) {
	var data dao.Article
	err := r.db.Get(ctx).First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *articleRepository) FindAll(ctx context.Context) ([]dao.Article, error) {
	var list []dao.Article
	err := r.db.Get(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *articleRepository) Paginate(ctx context.Context, limit, offset int) ([]dao.Article, int64, error) {
	var list []dao.Article
	var total int64

	tx := r.db.Get(ctx)
	if err := tx.Model(&dao.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := tx.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

func (r *articleRepository) Create(ctx context.Context, data *dao.Article) error {
	return r.db.Get(ctx).Create(data).Error
}

func (r *articleRepository) Update(ctx context.Context, data *dao.Article) error {
	return r.db.Get(ctx).Save(data).Error
}

func (r *articleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Get(ctx).Delete(&dao.Article{}, id).Error
}
