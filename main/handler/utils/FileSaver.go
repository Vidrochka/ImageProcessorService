package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

//FileSaver - file save manager
type FileSaver struct {
	logger *log.Logger
	config *Configuration
}

//CreateFileSaver - create file save manager
func CreateFileSaver(logger *log.Logger, config *Configuration) *FileSaver {
	saver := FileSaver{logger: logger, config: config}

	if _, err := os.Stat(config.FileSavePath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(config.FileSavePath, 0777); err != nil {
				saver.logger.Panicln(err)
				panic(err)
			}
		} else {
			saver.logger.Panicln(err)
			panic(err)
		}
	}

	saver.logger.Println("FileSaver created")

	return &saver
}

//SaveFile - save file data
func (saver *FileSaver) SaveFile(name, extension, data string) string {
	var err error

	hash := MD5(name + "." + extension + fmt.Sprint(time.Now()))
	saveFolder := saver.config.FileSavePath + hash + "/"

	if err = os.Mkdir(saveFolder, 0777); err != nil {
		saver.logger.Panicln(err)
		panic(err)
	}

	var file *os.File
	if file, err = os.Create(saveFolder + name + "." + extension); err != nil {
		saver.logger.Panicln(err)
		panic(err)
	}

	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		saver.logger.Panicln(err)
		panic(err)
	}

	return hash
}

//SavePreview - save previev image
func (saver *FileSaver) SavePreview(name, extension, data, hash string) {
	var err error

	saveFolder := saver.config.FileSavePath + hash + "/" + saver.config.PreviewFileFolder + "/"

	if _, err = os.Stat(saveFolder); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(saveFolder, 0777); err != nil {
				saver.logger.Panicln(err)
				panic(err)
			}
		} else {
			saver.logger.Panicln(err)
			panic(err)
		}
	}

	var file *os.File
	if _, err = os.Stat(saveFolder + name + "." + extension); err != nil {
		if os.IsNotExist(err) {
			saver.logger.Println("File not created yet. Creating...")

			if file, err = os.Create(saveFolder + name + "." + extension); err != nil {
				saver.logger.Panicln(err)
				panic(err)
			}
		} else {
			saver.logger.Panicln(err)
			panic(err)
		}
	} else {
		saver.logger.Println("File exist. Remove and create")

		if err = os.Remove(saveFolder + name + "." + extension); err != nil {
			saver.logger.Panicln(err)
			panic(err)
		}
		if file, err = os.Create(saveFolder + name + "." + extension); err != nil {
			saver.logger.Panicln(err)
			panic(err)
		}
	}

	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		saver.logger.Panicln(err)
		panic(err)
	}
}

//RestoreFile - restore file data
func (saver *FileSaver) RestoreFile(name, extension, hash string) string {
	var err error

	fileFolder := saver.config.FileSavePath + hash + "/"
	var file *os.File
	if file, err = os.Open(fileFolder + name + "." + extension); err != nil {
		saver.logger.Panicln(err)
		panic(err)
	}

	defer file.Close()

	var data []byte
	if data, err = ioutil.ReadAll(file); err != nil {
		saver.logger.Panicln(err)
		panic(err)
	}

	return string(data)
}

//MD5 - take MD5
func MD5(data string) string {
	h := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", h)
}
