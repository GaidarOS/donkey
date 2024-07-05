package main

import (
	"log/slog"
	"os"

	"receipt_store/config"
	"receipt_store/logger"
	"receipt_store/middleware"
	"receipt_store/routes"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/joho/godotenv/autoload"
)

var (
	slogger *slog.Logger
	Config  = swagger.Config{
		Next:     nil,
		BasePath: "/",
		FilePath: "./swagger.json",
		Path:     "docs",
		Title:    "Donkey API documentation",
		CacheAge: 3600, // Default to 1 hour
	}
)

func init() {
	// Load .env file

	// Make sure that the base directory exists
	if err := os.MkdirAll(config.AppConf.Dir, os.ModePerm); err != nil {
		slog.Error("Couldn't create the directory", err)
	}

	slogger = logger.Logger()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Donkey v1.0.0",

		// If run with the following param the list of routes will be printed in the log when starting the server
		// EnablePrintRoutes: true,
	})

	app.Use(recover.New())
	app.Use(swagger.New(Config))
	app.Use(csrf.New())
	app.Use(compress.New())
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "https://gofiber.io, https://gofiber.net",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))
	// pass the custom logger through the fiber context
	app.Use(slogfiber.New(slogger))

	api_v1 := app.Group("/api/v1")
	api_v1.Get("/download/*", middleware.TokenMiddleware, routes.DownloadFile)
	api_v1.Get("/list/*", middleware.TokenMiddleware, routes.ListFiles)
	api_v1.Post("/upload/*", middleware.TokenMiddleware, routes.SaveFile)
	api_v1.Delete("/delete/*", middleware.TokenMiddleware, routes.DeleteFiles)

	admin_v1 := api_v1.Group("/admin")
	admin_v1.Get("/", middleware.AdminMiddleware, routes.TokensList)
	admin_v1.Post("/", middleware.AdminMiddleware, routes.TokenCreate)
	admin_v1.Put("/", middleware.AdminMiddleware, routes.TokenEdit)
	admin_v1.Delete("/", middleware.AdminMiddleware, routes.TokenDelete)
	admin_v1.Get("/config", middleware.AdminMiddleware, routes.GetConfig)
	admin_v1.Post("/config", middleware.AdminMiddleware, routes.UpdateConfig)
	// Start server on port 3000
	slog.Debug("Starting the web-server!")
	err := app.Listen(":" + config.AppConf.Port)
	if err != nil {
		slog.Error("Couldn't start the fiber server", err)
	}

}
