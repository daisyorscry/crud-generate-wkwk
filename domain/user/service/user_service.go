package service

import (
	"context"

	"daisy/domain/user/dao"
	"daisy/domain/user/repository"
	"daisy/pkg/responses"
	"github.com/go-playground/validator/v10"
    "daisy/pkg/pkgErr"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	Paginate(ctx context.Context, limit, offset int) *responses.BaseResponse
	GetByID(ctx context.Context, id uint) *responses.BaseResponse
	Create(ctx context.Context, req *dao.CreateUserRequest) *responses.BaseResponse
	Update(ctx context.Context, req *dao.UpdateUserRequest) *responses.BaseResponse
	Delete(ctx context.Context, id uint) *responses.BaseResponse
}

type userService struct {
	repo      repository.UserRepository
	validator *validator.Validate
}

func NewUserService(repo repository.UserRepository,  validator *validator.Validate) UserService {
	return &userService{repo: repo, validator: validator}
}

func (s *userService) Paginate(ctx context.Context, limit, offset int) *responses.BaseResponse {
	items, total, err := s.repo.Paginate(ctx, limit, offset)
	if err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to fetch data", err)
	}
	return responses.ResponsePagination(fiber.StatusOK, "User list", items, total)
}

func (s *userService) GetByID(ctx context.Context, id uint) *responses.BaseResponse {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return responses.ResponseError(fiber.StatusNotFound, "User not found", err)
	}
	return responses.ResponseSuccess(fiber.StatusOK, "User detail", item)
}

func (s *userService) Create(ctx context.Context, req *dao.CreateUserRequest) *responses.BaseResponse {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return responses.ResponseValidation(
			fiber.ErrUnprocessableEntity.Code,
			"Validation failed",
			pkgErr.ParseValidationErrors(err),
		)
	}

	data := dao.ToUser(req)

	if err := s.repo.RunTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to create User", err)
	}

	return responses.ResponseSuccess(fiber.StatusCreated, "User created successfully", nil)
}

func (s *userService) Update(ctx context.Context, req *dao.UpdateUserRequest) *responses.BaseResponse {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return responses.ResponseValidation(
			fiber.ErrUnprocessableEntity.Code,
			"Validation failed",
			pkgErr.ParseValidationErrors(err),
		)
	}


	data := dao.ToUpdatedUser(req)
	if err := s.repo.RunTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Update(txCtx, data); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to update User", err)
	}

	return responses.ResponseSuccess(fiber.StatusOK, "User updated successfully", nil)
}

func (s *userService) Delete(ctx context.Context, id uint) *responses.BaseResponse {
	if err := s.repo.Delete(ctx, id); err != nil {
		return responses.ResponseError(fiber.StatusInternalServerError, "Failed to delete User", err)
	}
	return responses.ResponseSuccess(fiber.StatusOK, "User deleted successfully", nil)
}
