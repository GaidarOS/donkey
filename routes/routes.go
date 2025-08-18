package routes

import (
	"donkey/config"
	"donkey/database"
	"donkey/helper"
	thumb "donkey/thumbnails"
	"fmt"
	"log/slog"
	"os"
	path "path"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Welcome Donkey!",
	})
}

func UsersList(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "List of user tokens",
		"data":    config.AppConf.Users,
	})
}

func DownloadFile(c *fiber.Ctx) error {
	// Parse JSON request body
	param := c.Params("*")
	possible_thumbnails := ""
	if len(strings.Split(param, "/")) > 0 {
		possible_thumbnails = strings.Split(param, "/")[0]
	}
	if possible_thumbnails == "thumbnails" {
		return c.SendFile(path.Join(config.AppConf.Dir, "thumbnails", path.Base(param)))
	} else {
		// Use SendFile to send the specified file for download
		return c.SendFile(path.Join(config.AppConf.Dir, c.Params("*")))
	}
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
			thumb_path := path.Join(config.AppConf.Dir, "thumbnails", path.Base(filename))
			if strings.Contains(filename, ".pdf") {
				thumb_filename := strings.Replace(thumb_path, ".pdf", ".jpg", 1)
				err := helper.DeleteFile(thumb_filename)
				if err != nil {
					slog.Error("Failed to delete the thumbnail pdf", slog.String("file", filename), slog.Any("err", err))
					return
				}
			} else {
				err := helper.DeleteFile(thumb_path)
				if err != nil {
					slog.Error("Failed to delete the thumbnail", slog.String("file", filename), slog.Any("err", err))
					return
				}
			}
		}(path.Join(config.AppConf.Dir, c.Params("*"), filename))
	}

	return c.JSON(fiber.Map{"message": "Files deleted successfully"})
}

// Endpoint to list all files in a directoryEndpoint to list all files in a directory
func ListFiles(c *fiber.Ctx) error {

	// Read all files in the directory
	files, err := os.ReadDir(path.Join(config.AppConf.Dir, c.Params("*")))
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
	var thumbnails []string
	//TODO: this probably needs to become recursive to handle files in subfolders
	for _, file := range files {
		if file.IsDir() {
			folderNames = append(folderNames, file.Name())
			continue
		}
		fileNames = append(fileNames, path.Join(c.Params("*"), file.Name()))
		thumbname := file.Name()
		if strings.Contains(file.Name(), ".pdf") {
			thumbname = strings.Replace(file.Name(), ".pdf", ".jpg", 1)
		}
		thumbnails = append(thumbnails, path.Join("thumbnails", thumbname))
	}
	// Return the list of file names
	return c.JSON(fiber.Map{"files": fileNames, "folders": folderNames, "thumbnails": thumbnails, "path": path.Join(config.AppConf.Dir, c.Params("*"))})
}

func SaveFile(c *fiber.Ctx) error {

	// Parse the multipart form:
	form, err := c.MultipartForm()
	if err != nil {
		slog.Error("Error while getting files!", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Encountered an error while uploading files"})

	}
	// Get all files from "documents" key:
	files := form.File["file"]

	// Loop through files:
	for _, file := range files {
		if _, ok := config.AppConf.AllowedHeaderTypes[file.Header["Content-Type"][0]]; !ok {
			slog.Error("Bad Request")
			continue
		}
		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

		// Save the files to disk:
		save_path := path.Join(config.AppConf.Dir, c.Params("*"), file.Filename)
		if err := c.SaveFile(file, save_path); err != nil {
			slog.Error("Couldn't save files!", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to save file"})
		}
		if !strings.Contains(file.Filename, ".pdf") {
			err = thumb.GenerateThumbnailFromImage(save_path, "thumbnails")
			if err != nil {
				slog.Error("Couldn't create thumbnail!", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create thumbnail"})
			}
		} else {
			err = thumb.GenerateThumbnailFromPdf(save_path, "thumbnails")
			if err != nil {
				slog.Error("Couldn't create thumbnail from pdf!", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create thumbnail from pdf"})
			}
		}
	}
	return c.JSON(fiber.Map{"message": "File uploaded successfully"})
}

func Index(c *fiber.Ctx) error {
	go func(folder string) {

		files, err := os.ReadDir(path.Join(config.AppConf.Dir, folder))

		if err != nil {
			slog.Error("Couldn't create thumbnail!", err)
			// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create thumbnail"})
		}

		thumbs, err := os.ReadDir(path.Join(config.AppConf.Dir, "thumbnails"))

		if err != nil {
			slog.Error("Couldn't create thumbnail!", err)
			// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create thumbnail"})
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			for _, thumb := range thumbs {
				if strings.Split(thumb.Name(), ".")[0] == strings.Split(file.Name(), ".")[0] {
					continue
				}
			}

			if strings.Contains(file.Name(), ".pdf") {
				thumb.GenerateThumbnailFromPdf(path.Join(config.AppConf.Dir, folder, file.Name()), "thumbnails")
			} else {
				thumb.GenerateThumbnailFromImage(path.Join(config.AppConf.Dir, folder, file.Name()), "thumbnails")
			}
		}
	}(c.Params("*"))

	return c.JSON(fiber.Map{"message": "Indexing finished"})
}

func GetBooks(c *fiber.Ctx) error {

	var books []config.Book
	result := database.DB.Find(&books)
	slog.Info("This is the list of books found: %s", slog.Any("err", result.Error))

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch items",
			"error":   result.Error.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Found items",
		"items":   books,
	})

}

func CreateBook(c *fiber.Ctx) error {

	var req config.Book
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	book := config.Book{
		Title:  req.Title,
		Author: req.Author,
	}

	res := database.DB.Create(&book)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Book has been added"})
}
