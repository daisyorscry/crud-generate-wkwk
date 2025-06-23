package dao

type TodoResponse struct {
	ID uint `json:"id"`
	Id string `json:"id"`
	Title string `json:"title"`
	Is_done bool `json:"is_done"`
}
