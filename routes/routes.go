package routes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"receipt_store/config"
	"receipt_store/helper"

	"github.com/gofiber/fiber/v2"
)

type downloadRequest struct {
	Filename string `json:"filename"`
}

func DownloadFile(c *fiber.Ctx) error {
	// Parse JSON request body
	var request downloadRequest
	if err := c.BodyParser(&request); err != nil {
		slog.Error("Error parsing the filename", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Use SendFile to send the specified file for download
	return c.SendFile(fmt.Sprintf("%s/%s/%s", config.AppConf.Dir, c.Params("*"), request.Filename))
}

func DeleteFiles(c *fiber.Ctx) error {

	// Try to parse JSON request body
	var request struct {
		Files []string `json:"files"`
	}
	// Parse the body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}
	// Delete each file
	for _, filename := range request.Files {
		go func(filename string) {
			err := helper.DeleteFile(filename)
			if err != nil {
				slog.Error("Failed to delete the file", slog.String("file", filename), slog.Any("err", err))
				return
			}
		}(fmt.Sprintf("%s/%s/%s", config.AppConf.Dir, c.Params("*"), filename))
	}

	return c.JSON(fiber.Map{"message": "Files deleted successfully"})
}

// Endpoint to list all files in a directoryEndpoint to list all files in a directory
func ListFiles(c *fiber.Ctx) error {

	// Read all files in the directory
	files, err := os.ReadDir(fmt.Sprintf("%s/%s", config.AppConf.Dir, c.Params("*")))
	if err != nil {
		slog.Error("Couldn't read directory. Does it exist?", err)
		// error out if the directory doesn't exist
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"message": "Internal Server Error. Directory doesn't exist.",
				"error":   err.Error(),
			})
	}
	// Make a list of all the file names
	var fileNames []string
	var folderNames []string
	for _, file := range files {
		if file.IsDir() {
			folderNames = append(folderNames, file.Name())
			continue
		}
		fileNames = append(fileNames, file.Name())
	}
	// Return the list of file names
	return c.JSON(fiber.Map{"files": fileNames, "folders": folderNames, "path": fmt.Sprintf("%s/%s", config.AppConf.Dir, c.Params("*"))})
}

func SaveFile(c *fiber.Ctx) error {

	// Parse the multipart form:
	form, err := c.MultipartForm()
	if err != nil {
		slog.Error("Error while getting files!", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Encountered an error while uploading files"})

	}
	// Get all files from "documents" key:
	files := form.File["receipt"]

	// Loop through files:
	for _, file := range files {
		if _, ok := config.AppConf.AllowedHeaderTypes[file.Header["Content-Type"][0]]; !ok {
			slog.Error("Bad Request")
			continue
		}
		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

		// Save the files to disk:
		if err := c.SaveFile(file, fmt.Sprintf("%s/%s/%s", config.AppConf.Dir, c.Params("*"), file.Filename)); err != nil {
			slog.Error("Couldn't save files!", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
		}
	}
	return c.JSON(fiber.Map{"message": "File uploaded successfully"})
}

func UpdateConfig(c *fiber.Ctx) error {

	var conf config.Config

	// Parse the body
	if err := c.BodyParser(&conf); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}
	rq, _ := json.Marshal(conf)
	if err := helper.WriteToFile(config.AppConf.ConfFile, rq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	return c.JSON(fiber.Map{"message": "Configuration updated successfully"})
}

func GetConfig(c *fiber.Ctx) error {

	return c.JSON(config.AppConf)
}