package helper

import (
	"os"
	"receipt_store/logger"
)

var (
	slogger = logger.Logger()
)

func DeleteFile(fileName string) error {
	// Check if the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		slogger.Error("File not found")
		return err
	}
	err := os.Remove(fileName)
	return err
}

func WriteToFile(fileName string, data []byte) error {
	os.Remove(fileName)
	// Create a file if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		slogger.Error("Error creating the file", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		slogger.Error("Error writing to file!", err)
		return err
	}
	return nil
}
