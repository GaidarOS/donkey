package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

var (
	// Get the logLevel from env variables.
	// If nothing is set it will use the info level
	level = os.Getenv("LOG_LEVEL")
	opts  = slog.HandlerOptions{
		Level:     logLevel(level),
		AddSource: checkLoglevel(level),
		// The bellow changes the logger time to unix time for manipulation
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			// Match the key we want
			if attr.Key == slog.TimeKey {
				attr.Key = "date"                               // Rename time into date
				attr.Value = slog.Int64Value(time.Now().Unix()) // Set it to a int64 unix time
			}
			return attr
		},
	}
)

func logLevel(logLevel string) slog.Level {
	// Map the available levels
	levels := map[string]slog.Level{"ERROR": slog.LevelError, "DEBUG": slog.LevelDebug, "WARN": slog.LevelWarn, "":slog.LevelInfo, "INFO":slog.LevelInfo}
	return levels[strings.ToUpper(logLevel)]
}

func checkLoglevel(logLevel string) bool {
	// Check if level is set to debug
	// If true enable the addSource on the logs  
	if strings.ToUpper(logLevel) == "DEBUG" {
		return true
	}
	return false
}

func Logger() *slog.Logger {
	// Create the logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
	return logger
}
