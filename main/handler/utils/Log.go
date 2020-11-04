package utils

import (
	"log"
	"os"
)

//CreateLog - create logger
func CreateLog(logFilePath string) *log.Logger {
	var file *os.File
	var err error

	if _, err = os.Stat(logFilePath); err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(logFilePath)

			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		file, err = os.OpenFile(logFilePath, os.O_RDWR, 0666)
		file.Seek(0, 2)

		if err != nil {
			panic(err)
		}
	}

	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return logger
}
