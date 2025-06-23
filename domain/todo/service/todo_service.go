package service

import (
	"context"

	"daisy/domain/todo/dao"
	"daisy/domain/todo/repository"
	"daisy/pkg/responses"
	"github.com/go-playground/validator/v10"
    "daisy/pkg/pkgErr"

	"github.com/gofiber/fiber/v2"
)

type TodoService interface {
	Paginate(ctx context.Context, limit, offset int) *responses.BaseResponse
	GetByID(ctx context.Context, id uint) *responses.BaseResponse
	Create(ctx context.Context, req *dao.CreateTodoRequest) *responses.BaseResponse
	Update(ctx context.Context, req *dao.UpdateTodoRequest) *responses.BaseResponse
	Delete(ctx context.Context, id uint) *responses.BaseResponse
}

type todoService struct {
	repo      repository.TodoRepository
	validator *validator.Validate
}

func NewTodoService(repo repository.TodoRepository,  validator *validator.Validate) TodoService {
	return &todoService{repo: repo, validator: validator}
}

func (s *todoService) Paginate(ctx context.Context, limit, offset int) *responses.BaseResponse {
	items, total, err := s.repo.Paginate(ctx, limit, offset)
	if err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to fetch data", err)
	}
	return responses.ResponsePagination(fiber.StatusOK, "Todo list", items, total)
}

func (s *todoService) GetByID(ctx context.Context, id uint) *responses.BaseResponse {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return responses.ResponseError(fiber.StatusNotFound, "Todo not found", err)
	}
	return responses.ResponseSuccess(fiber.StatusOK, "Todo detail", item)
}

func (s *todoService) Create(ctx context.Context, req *dao.CreateTodoRequest) *responses.BaseResponse {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return responses.ResponseValidation(
			fiber.ErrUnprocessableEntity.Code,
			"Validation failed",
			pkgErr.ParseValidationErrors(err),
		)
	}

	data := dao.ToTodo(req)

	if err := s.repo.RunTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to create Todo", err)
	}

	return responses.ResponseSuccess(fiber.StatusCreated, "Todo created successfully", nil)
}

func (s *todoService) Update(ctx context.Context, req *dao.UpdateTodoRequest) *responses.BaseResponse {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return responses.ResponseValidation(
			fiber.ErrUnprocessableEntity.Code,
			"Validation failed",
			pkgErr.ParseValidationErrors(err),
		)
	}


	data := dao.ToUpdatedTodo(req)
	if err := s.repo.RunTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Update(txCtx, data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to update Todo", err)
	}

	return responses.ResponseSuccess(fiber.StatusOK, "Todo updated successfully", nil)
}

func (s *todoService) Delete(ctx context.Context, id uint) *responses.BaseResponse {
	if err := s.repo.Delete(ctx, id); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to delete Todo", err)
	}
	return responses.ResponseSuccess(fiber.StatusOK, "Todo deleted successfully", nil)
}
