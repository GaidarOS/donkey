package config

import (
	"donkey/helper"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {

	// Assuming config is an instance of your Config struct
	cfg := Config{
		Dir:                "uploads",
		Port:               "8080",
		Depth:              "3",
		ConfFile:           "../config-test.json",
		AllowedHeaderTypes: map[string]bool{"image/vnd.djvu": true},
	}
	cfg.WriteToConf()

	cfg.updateFromFile()
	testCases := map[string]struct {
		got      any
		expected any
		pass     bool
	}{
		"Directory": {
			got:      AppConf.Dir,
			expected: "uploads",
			pass:     true,
		},
		"Port": {
			got:      cfg.Port,
			expected: "8080",
			pass:     true,
		},
		"Depth": {
			got:      cfg.Depth,
			expected: "3",
			pass:     true,
		},
		"ConfFile": {
			got:      cfg.ConfFile,
			expected: "../config-test.json",
			pass:     true,
		},
		"AllowedHeaderTypes": {
			got:      cfg.AllowedHeaderTypes["image/vnd.djvu"],
			expected: true,
			pass:     true,
		},
		"NotAllowedHeaderTypes": {
			got:      cfg.AllowedHeaderTypes["audio/mp3"],
			expected: true,
			pass:     false,
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tC.pass, (tC.got == tC.expected))
		})
	}
	helper.DeleteFile(cfg.ConfFile)
}
