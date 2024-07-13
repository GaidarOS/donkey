package main

import (
	"donkey/config"
	"donkey/logger"
	"donkey/middleware"
	"donkey/routes"
	"log/slog"
	"os"
	"path"

	maps "golang.org/x/exp/maps"

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
	if err := os.MkdirAll(path.Join(config.AppConf.Dir, "thumbnails"), os.ModePerm); err != nil {
		slog.Error("Couldn't create thumbnails directory", err)
	}

	for _, user := range config.AppConf.Users {
		for _, folder := range maps.Keys(user.AccessPaths) {
			slog.Info("Creting user folder " + folder + " if it doesn't exist")
			if err := os.MkdirAll(path.Join(config.AppConf.Dir, folder), os.ModePerm); err != nil {
				slog.Error("Couldn't create thumbnails directory", err)
			}
		}
	}

	slogger = logger.Logger()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:   "Donkey v1.0.0",
		BodyLimit: 50 * 1024 * 1024,
		// StreamRequestBody:            true,
		// If run with the following param the list of routes will be printed in the log when starting the server
		// EnablePrintRoutes: true,
	})

	app.Use(recover.New())
	// app.Use(csrf.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://gofiber.io, http://127.0.0.1:5173/",
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
	api_v1.Post("/index/*", middleware.TokenMiddleware, routes.Index)

	admin_v1 := api_v1.Group("/admin")
	admin_v1.Get("/users", middleware.AdminMiddleware, routes.UsersList)
	admin_v1.Post("/user", middleware.AdminMiddleware, routes.UserCreate)
	admin_v1.Put("/user", middleware.AdminMiddleware, routes.UserUpdate)
	admin_v1.Delete("/user", middleware.AdminMiddleware, routes.UserDelete)
	admin_v1.Get("/config", middleware.AdminMiddleware, routes.GetConfig)
	admin_v1.Post("/config", middleware.AdminMiddleware, routes.UpdateConfig)

	slog.Debug("Starting the web-server!")
	err := app.Listen(":" + config.AppConf.Port)
	if err != nil {
		slog.Error("Couldn't start the fiber server", err)
	}
}
