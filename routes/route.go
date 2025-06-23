package routes

import (
	"daisy/domain"
	"daisy/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, conn database.Connection) {
	d := domain.NewDomain(conn)
	api := app.Group("/api")
	// todo routes
	todoGroup := api.Group("/todo")
	todoGroup.Post("/", d.TodoHandler.Create)
	todoGroup.Put("/:id", d.TodoHandler.Update)
	todoGroup.Get("/:id", d.TodoHandler.GetByID)
	todoGroup.Get("/", d.TodoHandler.Paginate)
	todoGroup.Delete("/:id", d.TodoHandler.Delete)

	// user routes
	apiGroup := api.Group("/user")
	apiGroup.Post("/", d.UserHandler.Create)
	apiGroup.Put("/:id", d.UserHandler.Update)
	apiGroup.Get("/:id", d.UserHandler.GetByID)
	apiGroup.Get("/", d.UserHandler.Paginate)
	apiGroup.Delete("/:id", d.UserHandler.Delete)

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
