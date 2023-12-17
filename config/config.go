package config

import (
	"log"
	"os"
)

type Config struct {
	Dir                string
	Port               string
	Depth              string
	ConfFile           string
	AllowedHeaderTypes map[string]bool
}

var (
	AppConf            Config
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
	AppConf = Config{
		Dir:  os.Getenv("DIR"),
		Port: os.Getenv("PORT"),
	}
	log.Println(AppConf)

}
