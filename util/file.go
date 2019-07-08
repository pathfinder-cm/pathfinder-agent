package util

import (
	"io"
	"os"
)

func WriteToFile(filename string, data string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return file, err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return file, err
	}
	return file, nil
}
