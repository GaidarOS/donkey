package config

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"receipt_store/helper"
	"receipt_store/logger"

	"github.com/fsnotify/fsnotify"
)

type Config struct {
	Dir                string          `json:"dir" default:"uploads"`
	Port               string          `json:"port" default:"8080"`
	Depth              string          `json:"depth" default:"3"`
	ConfFile           string          `json:"confFile" default:"./config.json"`
	AllowedHeaderTypes map[string]bool `json:"allowedHeaderTypes"`
}

var (
	slogger             = logger.Logger()
	default_config_path = "./config.json"

	AppConf = Config{
		ConfFile: default_config_path,
	}

	allowedHeaderTypes = map[string]bool{
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
	filesExist, err := initChecks()
	if err != nil {
		slogger.Error("Critical Error!", err)
		os.Exit(1)
	}

	if filesExist {
		AppConf.updateFromFile()
		if len(AppConf.AllowedHeaderTypes) == 0 {
			AppConf.AllowedHeaderTypes = allowedHeaderTypes
		}
		slogger.Debug("Config", slog.Any("conf", AppConf))
		AppConf.watchConfig()
	} else {
		slogger.Error("No conifg files found. Will continue with defaults")
	}
}

func initChecks() (bool, error) {
	// Checks if the file or the environment variables exist
	if _, err := os.Stat(AppConf.ConfFile); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		_, present := os.LookupEnv("PORT")
		if present {
			return false, nil
		}
		return false, errors.New("no port variable found")
	} else {
		return false, errors.New("neither config file, not env variables present")
	}
}

func (c *Config) updateFromFile() {

	file, err := os.ReadFile(c.ConfFile)
	if err != nil {
		slogger.Error("No file: ", err)
	}

	err = json.Unmarshal([]byte(file), &c)
	if err != nil {
		slogger.Error("Could not unmarshal the file", err)
	}

}

func (c *Config) watchConfig() {

	slogger.Info("Will create the config file if changed")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slogger.Warn("Error creating file watcher", err)
	}
	defer watcher.Close()

	// check for updates on the configFile and update the config when detected
	go func() {
		wtchr, err := fsnotify.NewWatcher()
		if err != nil {
			slogger.Error("Error creating file watcher", err)
		}
		addToWatcher(wtchr, c.ConfFile)
		defer wtchr.Close()
		for {
			select {
			// watch for events
			case event := <-wtchr.Events:
				if event.Op.String() == "CHMOD" {
					slogger.Info("config.json changes detected", slog.Any("event", event.Op.String()))
					c.updateFromFile()
				}
			case err := <-wtchr.Errors:
				slogger.Error("Config didn't change changed", err)
			}
			addToWatcher(wtchr, c.ConfFile)
		}
	}()

}

func addToWatcher(watcher *fsnotify.Watcher, filename string) {
	if err := watcher.Add(filename); err != nil {
		slogger.Error("Could not add file to the watcher", err)
	}
}

func (c *Config) writeToConf() error {
	// Marshal the struct to JSON
	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		slogger.Error("Error marshaling JSON:", err)
		return err
	}

	helper.WriteToFile(c.ConfFile, jsonData)

	slogger.Info("Config successfully saved to", slog.Any("filename", c.ConfFile))
	return nil
}
