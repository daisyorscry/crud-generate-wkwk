package dao

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email string `json:"email" validate:"required"`
	Age int `json:"age" validate:"required"`
	Is_active bool `json:"is_active" validate:"required"`
}

type UpdateUserRequest struct {
	ID uint `json:"id"`
	Username string `json:"username" validate:"required"`
	Email string `json:"email" validate:"required"`
	Age int `json:"age" validate:"required"`
	Is_active bool `json:"is_active" validate:"required"`
}
