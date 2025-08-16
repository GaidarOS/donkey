package config

import (
	"donkey/logger"
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
)

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
		slogger.Error("Critical Error!", slog.Any("err", err))
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
		slogger.Error("No config files found. Will continue with defaults")
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

	slogger.Debug("Updating configuration from file")
	file, err := os.ReadFile(c.ConfFile)
	if err != nil {
		slogger.Error("No file: ", slog.Any("err", err))
	}

	err = json.Unmarshal([]byte(file), &c)
	if err != nil {
		slogger.Error("Could not unmarshal the file", slog.Any("err", err))
	}
	// We can't seriously print the config as standard info
	// when user information is contained in it.
	slogger.Debug("updated config", slog.Any("config", c))
}

func (c *Config) watchConfig() {

	slogger.Info("Will create the config file if changed")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slogger.Warn("Error creating file watcher", slog.Any("err", err))
	}
	defer watcher.Close()

	// check for updates on the configFile and update the config when detected
	go func() {
		wtchr, err := fsnotify.NewWatcher()
		if err != nil {
			slogger.Error("Error creating file watcher", slog.Any("err", err))
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
				slogger.Error("Config didn't change changed", slog.Any("err", err))
			}
			addToWatcher(wtchr, c.ConfFile)
		}
	}()
}

func addToWatcher(watcher *fsnotify.Watcher, filename string) {
	if err := watcher.Add(filename); err != nil {
		slogger.Error("Could not add file to the watcher", slog.Any("err", err))
	}
}
