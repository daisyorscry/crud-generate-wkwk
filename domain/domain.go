package domain

import (
	userRest "daisy/domain/user/handler/rest"
	userRepo "daisy/domain/user/repository"
	userService "daisy/domain/user/service"
	"daisy/pkg/database"

	todoRest "daisy/domain/todo/handler/rest"
	todoRepo "daisy/domain/todo/repository"
	todoService "daisy/domain/todo/service"

	"github.com/go-playground/validator/v10"
)

type Domain struct {
	UserHandler *userRest.UserHandler
	TodoHandler *todoRest.TodoHandler
}

func NewDomain(conn database.Connection) *Domain {
	d := &Domain{}

	userRepo := userRepo.New(conn)
	userService := userService.NewUserService(userRepo, validator.New())
	userHandler := userRest.NewUserHandler(userService)
	d.UserHandler = userHandler

	todoRepo := todoRepo.New(conn)
	todoService := todoService.NewTodoService(todoRepo, validator.New())
	todoHandler := todoRest.NewTodoHandler(todoService)
	d.TodoHandler = todoHandler
	return d
}
