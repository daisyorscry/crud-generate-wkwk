package dao

type CreateArticleRequest struct {
	Title string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Is_active bool `json:"is_active" validate:"required"`
}

type UpdateArticleRequest struct {
	ID uint `json:"id"`
	Title string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Is_active bool `json:"is_active" validate:"required"`
}
