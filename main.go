package main

import (
	"daisy/pkg/generator"
	"fmt"
	"net/http"
)

func main() {
	// cfg := config.Get()
	// conn, err := database.NewConnection(cfg.Database)
	// if err != nil {
	// 	log.Fatalf("failed to connect to database: %v", err)
	// }
	// defer func() {
	// 	sqlDB, err := conn.Get(context.Background()).DB()
	// 	if err == nil {
	// 		if closeErr := sqlDB.Close(); closeErr != nil {
	// 			slog.Error("Failed to close DB", "err", closeErr)
	// 		} else {
	// 			slog.Info("Database connection closed")
	// 		}
	// 	}
	// }()

	// app := fiber.New()
	// routes.Setup(app, conn)
	http.HandleFunc("/generate-model", generator.HandleGenerateModel)
	http.HandleFunc("/generate-repository", generator.HandleGenerateRepository)
	http.HandleFunc("/generate-all", generator.HandleGenerateAll)

	fmt.Println("🚀 Running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
