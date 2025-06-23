package dao

type UserResponse struct {
	ID uint `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Age int `json:"age"`
	Is_active bool `json:"is_active"`
}
