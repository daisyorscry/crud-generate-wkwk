package dao

func ToUser(req *CreateUserRequest) *User {
	return &User{
		Username: req.Username,
		Email: req.Email,
		Age: req.Age,
		Is_active: req.Is_active,
	}
}

func ToUpdatedUser(req *UpdateUserRequest) *User {
	return &User{
		ID: req.ID,
		Username: req.Username,
		Email: req.Email,
		Age: req.Age,
		Is_active: req.Is_active,
	}
}
