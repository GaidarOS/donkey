package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {

	// Assuming config is an instance of your Config struct
	cfg := Config{
		Dir:                "uploads",
		Port:               "8080",
		Depth:              "3",
		ConfFile:           "config.json",
		AllowedHeaderTypes: map[string]bool{"image/vnd.djvu":  true},
	}
	fileneame := writeToFile(cfg)

	AppConf.updateFromFile(fileneame)
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
			got:      AppConf.Port,
			expected: "8080",
			pass:     true,
		},
		"Depth": {
			got:      AppConf.Depth,
			expected: "3",
			pass:     true,
		},
		"ConfFile": {
			got:      AppConf.ConfFile,
			expected: "config.json",
			pass:     true,
		},
		"AllowedHeaderTypes": {
			got:      AppConf.AllowedHeaderTypes["image/vnd.djvu"],
			expected: true,
			pass:     true,
		},
		"NotAllowedHeaderTypes": {
			got:      AppConf.AllowedHeaderTypes["audio/mp3"],
			expected: true,
			pass:     false,
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tC.pass, (tC.got == tC.expected))
		})
	}
	deleteFile(fileneame)
}
