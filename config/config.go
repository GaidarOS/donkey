package config

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
)

type Config struct {
	Dir                string
	Port               string
	Depth              string
	ConfFile           string
	AllowedHeaderTypes map[string]bool
}

var (
	default_config_path = "./config.json"

	AppConf = Config{
		ConfFile: default_config_path,
	}

	AllowedHeaderTypes = map[string]bool{
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
)

func init() {
	// Default config path
	// check if file exists
	f, err := initChecks()
	if err != nil {
		slog.Error("Critical Error!", err)
		os.Exit(1)
	}

	if f {
		AppConf.updateFromFile(AppConf.ConfFile)
		if len(AppConf.AllowedHeaderTypes) == 0 {
			AppConf.AllowedHeaderTypes = AllowedHeaderTypes
		}
		slog.Info("Config", slog.Any("conf", AppConf))
		AppConf.watchConfig()
	}else{

	}
}

func initChecks() (bool, error) {
	// Checks if the file or the environment variables exist
	if _, err := os.Stat(default_config_path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist){
		_, present := os.LookupEnv("PORT")
		if present {
			return false, nil
		}
		return false, errors.New("No PORT variable found")
	} else {
		return false, errors.New("Neither config file, not env variables present")
	}
}

func (c *Config) updateFromFile(filename string) {

	file, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("No file: ", err)
	}

	err = json.Unmarshal([]byte(file), &c)
	if err != nil {
		slog.Error("Err:", err)
	}

}

func (c *Config) watchConfig() {

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
		addToWatcher(wtchr, c.ConfFile)
		defer wtchr.Close()
		for {
			select {
			// watch for events
			case event := <-wtchr.Events:
				if event.Op.String() == "CHMOD" {
					log.Print("Webapp.yaml changes detected", event.Op.String())
					c.updateFromFile(c.ConfFile)
				}
			case err := <-wtchr.Errors:
				log.Print("webapp didn't change changed", err)
			}
			addToWatcher(wtchr, c.ConfFile)
		}
	}()
	// // Make sure the job never ends
	select {}
}

func addToWatcher(watcher *fsnotify.Watcher, filename string) {
	if err := watcher.Add(filename); err != nil {
		log.Println("Could not add file to the watcher", err)
	}
}

func writeToFile(config Config) string {
	// Marshal the struct to JSON
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		slog.Error("Error marshaling JSON:", err)
		return ""
	}

	// Write JSON data to a file
	fileName := "config-test.json"
	// err = writeToFile(fileName, jsonData)
	file, err := os.Create(fileName)
	if err != nil {
		slog.Error("Error creating the file", err)
		return ""
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		slog.Error("Error writting to file!", err)
		return fileName
	}

	slog.Info("Config successfully saved to", slog.Any("filename", fileName))
	return fileName
}

func deleteFile(fileName string) error {
	err := os.Remove(fileName)
	return err
}
