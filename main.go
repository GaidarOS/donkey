package main

import (
	"donkey/config"
	"donkey/database"
	"donkey/logger"
	"donkey/middleware"
	"donkey/routes"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/gofiber/fiber/v2/middleware/recover"
	// Automatically load .env file
	_ "github.com/joho/godotenv/autoload"
)

var (
	slogger *slog.Logger
)

func init() {

	slogger = logger.Logger()
	// Initialize database. Based on variables in the .env file
	database.Connect()

}

func main() {
	api := fiber.New(fiber.Config{
		AppName:   "Donkey v1.0.0",
		BodyLimit: 50 * 1024 * 1024,
		// Disable the default startup message
		DisableStartupMessage: true,
		// StreamRequestBody:            true,
		// If run with the following param the list of routes will be printed in the log when starting the server
		// EnablePrintRoutes: true,
	})

	api.Use(recover.New())
	// api.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Lax",
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.UUIDv4,
	// }))
	api.Use(compress.New())
	api.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173/",
		AllowHeaders:     "Cache-Control, Accept-Encoding, Origin, Content-Type, Accept, Token, Content-Length, Accept-Encoding, X-Csrf-Token, Authorization",
		AllowCredentials: true,
		AllowMethods:     "GET, HEAD, PUT, PATCH, POST, DELETE",
		MaxAge:           0,
	}))

	// OTEL compatible logs
	logConfig := slogfiber.Config{
		WithSpanID:       true,
		WithTraceID:      true,
		ServerErrorLevel: slog.LevelError,
	}

	// Verbose loggin on DEBUG
	if slog.LevelError == slog.LevelDebug {
		logConfig = slogfiber.Config{
			WithRequestBody:    true,
			WithResponseBody:   true,
			WithRequestHeader:  true,
			WithResponseHeader: true,
		}
	}

	// pass the custom logger through the fiber context
	api.Use(slogfiber.NewWithConfig(slogger, logConfig))

	// Auth
	api_v1 := api.Group("/api/v1")

	auth := api_v1.Group("/auth")
	auth.Post("/login", routes.Login)
	auth.Post("/register", routes.Register)

	api_v1.Get("/books", routes.GetBooks)

	api_v1.Use(middleware.JWTProtected)
	api_v1.Post("/books", middleware.JWTProtected, routes.CreateBook)
	api_v1.Get("/download/*", routes.DownloadFile)

	slog.Debug("Starting the web-server!")

	api.Hooks().OnListen(func(listenData fiber.ListenData) error {

		// ASCII generated with
		// https://www.asciiart.eu/text-to-ascii-art
		// Font: Ghost
		// http://www.jave.de/figlet/fonts/details/ghost.html

		fmt.Println("\n" +
			" _ .-') _                    .-') _ .-. .-')     ('-.\n" +
			"( (  OO) )                  ( OO ) )\\  ( OO )  _(  OO)\n" +
			" \\     .'_  .-'),-----. ,--./ ,--,' ,--. ,--. (,------. ,--.   ,--.\n" +
			" ,`'--..._)( OO'  .-.  '|   \\ |  |\\ |  .'   /  |  .---'  \\  `.'  /\n" +
			" |  |  \\  '/   |  | |  ||    \\|  | )|      /,  |  |    .-')     /\n" +
			" |  |   ' |\\_) |  |\\|  ||  .     |/ |     ' _)(|  '--.(OO  \\   /\n" +
			" |  |   / :  \\ |  | |  ||  |\\    |  |  .   \\   |  .--' |   /  /\n" +
			" |  '--'  /   `'  '-'  '|  | \\   |  |  |\\   \\  |  `---.`-./  /\n" +
			" `-------'      `-----' `--'  `--'  `--' '--'  `------'  `--'\n" +
			" ")
		slog.Info("To connect to the server use => http://" + listenData.Host + ":" + listenData.Port)
		return nil
	})

	err := api.Listen(":" + config.AppConf.Port)
	if err != nil {
		slog.Error("Couldn't start the fiber server", slog.Any("err", err))
	}
}
