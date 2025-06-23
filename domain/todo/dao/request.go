package dao

type CreateTodoRequest struct {
	Id string `json:"id" validate:"required"`
	Title string `json:"title" validate:"required"`
	Is_done bool `json:"is_done" validate:"required"`
}

type UpdateTodoRequest struct {
	ID uint `json:"id"`
	Id string `json:"id" validate:"required"`
	Title string `json:"title" validate:"required"`
	Is_done bool `json:"is_done" validate:"required"`
}
