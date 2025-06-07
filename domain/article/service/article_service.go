package service

import (
	"context"

	"daisy/domain/article/dao"
	"daisy/domain/article/repository"
	"daisy/pkg/responses"
	"github.com/go-playground/validator/v10"
    "daisy/pkg/pkgErr"

	"github.com/gofiber/fiber/v2"
)

type ArticleService interface {
	Paginate(ctx context.Context, limit, offset int) *responses.BaseResponse
	GetByID(ctx context.Context, id uint) *responses.BaseResponse
	Create(ctx context.Context, req *dao.CreateArticleRequest) *responses.BaseResponse
	Update(ctx context.Context, req *dao.UpdateArticleRequest) *responses.BaseResponse
	Delete(ctx context.Context, id uint) *responses.BaseResponse
}

type articleService struct {
	repo      repository.ArticleRepository
	validator *validator.Validate
}

func NewArticleService(repo repository.ArticleRepository,  validator *validator.Validate) ArticleService {
	return &articleService{repo: repo, validator: validator}
}

func (s *articleService) Paginate(ctx context.Context, limit, offset int) *responses.BaseResponse {
	items, total, err := s.repo.Paginate(ctx, limit, offset)
	if err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to fetch data", err)
	}
	return responses.ResponsePagination(fiber.StatusOK, "Article list", items, total)
}

func (s *articleService) GetByID(ctx context.Context, id uint) *responses.BaseResponse {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return responses.ResponseError(fiber.StatusNotFound, "Article not found", err)
	}
	return responses.ResponseSuccess(fiber.StatusOK, "Article detail", item)
}

func (s *articleService) Create(ctx context.Context, req *dao.CreateArticleRequest) *responses.BaseResponse {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return responses.ResponseValidation(
			fiber.ErrUnprocessableEntity.Code,
			"Validation failed",
			pkgErr.ParseValidationErrors(err),
		)
	}

	data := dao.ToArticle(req)

	if err := s.repo.RunTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to create Article", err)
	}

	return responses.ResponseSuccess(fiber.StatusCreated, "Article created successfully", nil)
}

func (s *articleService) Update(ctx context.Context, req *dao.UpdateArticleRequest) *responses.BaseResponse {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return responses.ResponseValidation(
			fiber.ErrUnprocessableEntity.Code,
			"Validation failed",
			pkgErr.ParseValidationErrors(err),
		)
	}


	data := dao.ToUpdatedArticle(req)
	if err := s.repo.RunTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Update(txCtx, data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to update Article", err)
	}

	return responses.ResponseSuccess(fiber.StatusOK, "Article updated successfully", nil)
}

func (s *articleService) Delete(ctx context.Context, id uint) *responses.BaseResponse {
	if err := s.repo.Delete(ctx, id); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to delete Article", err)
	}
	return responses.ResponseSuccess(fiber.StatusOK, "Article deleted successfully", nil)
}
