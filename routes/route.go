package routes

import (
	"daisy/domain/user/handler/rest"
	"daisy/domain/user/repository"
	"daisy/domain/user/service"
	"daisy/pkg/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, conn database.Connection) {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	userRepo := repository.New(conn)
	userService := service.New(userRepo, validator.New())
	userHandler := rest.NewUserHandler(userService)

	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
}
