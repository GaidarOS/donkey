package main

import (
	"log/slog"
	"os"

	"receipt_store/config"
	"receipt_store/logger"
	"receipt_store/middleware"
	"receipt_store/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/joho/godotenv/autoload"
)

var (
	slogger *slog.Logger
)

func init() {
	// Load .env file

	// Make sure that the base directory exists
	if err := os.MkdirAll(config.AppConf.Dir, os.ModePerm); err != nil {
		slog.Error("Couldn't create the directory", err)
	}

	// Make sure that the base directory exists
	if err := os.MkdirAll(config.AppConf.Dir+"/thumbnails", os.ModePerm); err != nil {
		slog.Error("Couldn't create thumbnails directory", err)
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
	// app.Use(csrf.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://gofiber.io, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Token, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		AllowCredentials: true,
		AllowMethods:     "GET, HEAD, PUT, PATCH, POST, DELETE",
		MaxAge:           0,
	}))

	// pass the custom logger through the fiber context
	app.Use(slogfiber.New(slogger))

	api_v1 := app.Group("/api/v1")
	api_v1.Post("/login", basicauth.New(basicauth.Config{
		Authorizer: func(username string, password string) bool {
			user, err := config.AppConf.FindStructByName(username)
			if err != nil {
				return false
			}
			return user.Password == password
		},
	}), routes.Login)
	api_v1.Get("/download/*", middleware.TokenMiddleware, routes.DownloadFile)
	api_v1.Get("/list/*", middleware.TokenMiddleware, routes.ListFiles)
	api_v1.Post("/upload/*", middleware.TokenMiddleware, routes.SaveFile)
	api_v1.Delete("/delete/*", middleware.TokenMiddleware, routes.DeleteFiles)
	api_v1.Get("/user/*", middleware.TokenMiddleware, routes.GetUser)

	admin_v1 := api_v1.Group("/admin")
	admin_v1.Get("/", middleware.AdminMiddleware, routes.UsersList)
	admin_v1.Post("/", middleware.AdminMiddleware, routes.UserCreate)
	admin_v1.Put("/", middleware.AdminMiddleware, routes.UserEdit)
	admin_v1.Delete("/", middleware.AdminMiddleware, routes.UserDelete)
	admin_v1.Get("/config", middleware.AdminMiddleware, routes.GetConfig)
	admin_v1.Post("/config", middleware.AdminMiddleware, routes.UpdateConfig)
	// Start server on port 3000
	slog.Debug("Starting the web-server!")
	err := app.Listen(":" + config.AppConf.Port)
	if err != nil {
		slog.Error("Couldn't start the fiber server", err)
	}
}
