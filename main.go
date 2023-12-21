package main

import (
	"log/slog"
	"os"
	"receipt_store/config"
	"receipt_store/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", err)
	}

	// Make sure that the base directory exists
	slog.Debug(os.Getenv("DIR"))
	if err := os.MkdirAll(config.AppConf.Dir, os.ModePerm); err != nil {
		slog.Error("Couldn't create the directory", err)
	}

}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Donkey",
	})

	app.Use(recover.New())

	api_v1 := app.Group("/api/v1")
	api_v1.Get("/download", routes.DownloadFile)

	// Route to handle file uploads
	api_v1.Post("/upload", routes.SaveFile)
	api_v1.Delete("/delete", routes.DeleteFile)
	api_v1.Post("/config", routes.UpdateConfig)
	api_v1.Get("/list", routes.ListFiles)

	// Start server on port 3000
	err := app.Listen(":" + config.AppConf.Port)
	if err != nil {
		slog.Error("Couldn't start the fiber server", err)
	}
}
