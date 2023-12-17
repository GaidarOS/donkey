package routes

import (
	"fmt"
	"log"
	"os"

	"receipt_store/config"

	"github.com/gofiber/fiber/v2"
)

func DownloadFile(c *fiber.Ctx) error {
	return nil
}

func DeleteFile(c *fiber.Ctx) error {
	return nil
}

// Endpoint to list all files in a directoryEndpoint to list all files in a directory
func ListFiles(c *fiber.Ctx) error {

	// Read all files in the directory
	files, err := os.ReadDir(config.AppConf.Dir)
	if err != nil {
		log.Println(err)
		// error out if the directory doesn't exist
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"message": "Internal Server Error. Directory doesn't exist.",
				"error":   err.Error(),
			})
	}
	// Make a list of all the file names
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	// Return the list of file names
	return c.JSON(fiber.Map{"files": fileNames})

}

func SaveFile(c *fiber.Ctx) error {

	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {

		// Get all files from "documents" key:
		files := form.File["receipt"]

		// Loop through files:
		for _, file := range files {
			if _, ok := config.AllowedHeaderTypes[file.Header["Content-Type"][0]]; !ok {
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

func UpdateConfig(c *fiber.Ctx) error  {

	return nil
}