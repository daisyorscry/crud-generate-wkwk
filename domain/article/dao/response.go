package dao

type ArticleResponse struct {
	ID uint `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Is_active bool `json:"is_active"`
}
