package routes

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func DownloadFile(c *fiber.Ctx) error {
	return nil
}

func SaveFile(c *fiber.Ctx) error {

	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {

		// Get all files from "documents" key:
		files := form.File["receipt"]

		allowedHeaderTypes := map[string]bool{
			"image/gif":       true,
			"image/jpeg":      true,
			"image/png":       true,
			"image/tiff":      true,
			"image/x-icon":    true,
			"image/vnd.djvu":  true,
			"image/svg+xml":   true,
			"image/jpg":       true,
			"application/pdf": true,
		}

		// Loop through files:
		for _, file := range files {
			if _, ok := allowedHeaderTypes[file.Header["Content-Type"][0]]; !ok {
				log.Println("Bad Request")
				continue
			}
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

			// Save the files to disk:
			if err := c.SaveFile(file, fmt.Sprintf("%s/%s", os.Getenv("DIR"), file.Filename)); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}
		}
		return err
	}

	return c.JSON(fiber.Map{"message": "File uploaded successfully"})
}
