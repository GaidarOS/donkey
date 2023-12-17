package routes

import (
	"fmt"
	"log"
	"os"

	"receipt_store/config"

	"github.com/fsnotify/fsnotify"
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

func UpdateConfig(c *fiber.Ctx) error {

	log.Println("Will create the config file if changed")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Print("Error creating file watcher", err)
	}
	defer watcher.Close()

	// check for updates on the webapp.yaml and update the webapp when detected
	go func() {
		wtchr, err := fsnotify.NewWatcher()
		if err != nil {
			log.Print("Error creating file watcher", err)
		}
		addToWatcher(wtchr, config.AppConf.ConfFile)
		defer wtchr.Close()
		for {
			select {
			// watch for events
			case event := <-wtchr.Events:
				if event.Op.String() == "CHMOD" {
					log.Print("Webapp.yaml changes detected", event.Op.String())
					reloadConfig()
				}
			case err := <-wtchr.Errors:
				log.Print("webapp didn't change changed", err)
			}
			addToWatcher(wtchr, config.AppConf.ConfFile)
		}
	}()
	// // Make sure the job never ends
	select {}
}

func addToWatcher(watcher *fsnotify.Watcher, filename string) {
	if err := watcher.Add(filename); err != nil {
		log.Printf("Could not add file to the watcher", err)
	}
}

func reloadConfig() {

}
