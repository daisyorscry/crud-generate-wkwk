package rest

import (
	"daisy/domain/todo/dao"
	"daisy/domain/todo/service"
	"daisy/pkg/responses"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	service service.TodoService
}

func NewTodoHandler(s service.TodoService) *TodoHandler {
	return &TodoHandler{service: s}
}

func (h *TodoHandler) Create(c *fiber.Ctx) error {
	var req dao.CreateTodoRequest

	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	resp := h.service.Create(c.Context(), &req)

	if resp.Errors != nil {
		return responses.ValidationError(c, resp.Code, resp.Message, resp.Errors)
	}
	if resp.Error != nil {
		return responses.Error(c, resp.Code, resp.Message, fmt.Errorf("%s", resp.Error.Detail))
	}

	return responses.Success(c, resp.Code, resp.Message, resp.Data)
}

func (h *TodoHandler) Update(c *fiber.Ctx) error {
	var req dao.UpdateTodoRequest

	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	resp := h.service.Update(c.Context(), &req)

	if resp.Errors != nil {
		return responses.ValidationError(c, resp.Code, resp.Message, resp.Errors)
	}
	if resp.Error != nil {
		return responses.Error(c, resp.Code, resp.Message, fmt.Errorf("%s", resp.Error.Detail))
	}

	return responses.Success(c, resp.Code, resp.Message, resp.Data)
}

func (h *TodoHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return responses.Error(c, fiber.StatusBadRequest, "Invalid ID param", err)
	}

	resp := h.service.GetByID(c.Context(), uint(id))

	if resp.Error != nil {
		return responses.Error(c, resp.Code, resp.Message, fmt.Errorf("%s", resp.Error.Detail))
	}

	return responses.Success(c, resp.Code, resp.Message, resp.Data)
}

func (h *TodoHandler) Paginate(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	resp := h.service.Paginate(c.Context(), limit, offset)

	if resp.Error != nil {
		return responses.Error(c, resp.Code, resp.Message, fmt.Errorf("%s", resp.Error.Detail))
	}

	return responses.Paginate(c, resp.Code, resp.Message, resp.Data, resp.Total)
}

func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return responses.Error(c, fiber.StatusBadRequest, "Invalid ID param", err)
	}

	resp := h.service.Delete(c.Context(), uint(id))

	if resp.Error != nil {
		return responses.Error(c, resp.Code, resp.Message, fmt.Errorf("%s", resp.Error.Detail))
	}

	return responses.Success(c, resp.Code, resp.Message, nil)
}
