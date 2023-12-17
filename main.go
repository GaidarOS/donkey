package main

import (
	"log"
	"os"
	"receipt_store/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Make sure that the base directory exists
	log.Println(os.Getenv("DIR"))
	if err := os.MkdirAll(os.Getenv("DIR"), os.ModePerm); err != nil {
		log.Fatal(err)
	}

}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Gaidaros",
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
	err := app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
