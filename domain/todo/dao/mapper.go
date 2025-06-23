package dao

func ToTodo(req *CreateTodoRequest) *Todo {
	return &Todo{
		Id: req.Id,
		Title: req.Title,
		Is_done: req.Is_done,
	}
}

func ToUpdatedTodo(req *UpdateTodoRequest) *Todo {
	return &Todo{
		ID: req.ID,
		Id: req.Id,
		Title: req.Title,
		Is_done: req.Is_done,
	}
}
