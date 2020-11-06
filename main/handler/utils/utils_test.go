package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/anthill-com/ImageProcessorService/main/handler/dto"
)

const _cfgPath string = "./test_config.toml"
const _logPath string = "./test.log"
const _dbPath string = "./test_db.db"
const _extensionsCollection string = "jpg,png,exe"

func TestConfig(t *testing.T) {
	var err error
	var file *os.File

	if file, err = os.Create(_cfgPath); err != nil {
		t.Error("Cant create file: " + err.Error())
	}

	defer os.Remove(_cfgPath)
	defer file.Close()

	file.WriteString(
		"\"LogFilePath\" = \"./service.log\"\n" +
			"\"Port\" = \"4183\"\n" +
			"\"ServedURL\" = \"/\"\n" +
			"\"ReadTimeout\" = 10\n" +
			"\"WriteTimeout\" = 10\n" +
			"\"FileSaveExtensionList\" = \"jpg,jpeg,png\"\n" +
			"\"ScaledImageRestoreExtension\" = \"jpg,jpeg,png\"\n" +
			"\"ScaledImageH\" = 100\n" +
			"\"ScapedImageW\" = 100\n")

	config := LoadConfiguration(_cfgPath)

	if config.LogFilePath != "./service.log" {
		t.Fatal("LogFilePath must be ./service.log but " + config.LogFilePath)
	}

	if config.Port != "4183" {
		t.Fatal("Port must be 4183 but " + config.Port)
	}

	if config.ScaledImageH != 100 {
		t.Fatal("ScaledImageH must be 100 but " + fmt.Sprint(config.ScaledImageH))
	}
}

func TestLog(t *testing.T) {
	var err error

	logger, logFile := CreateLog(_logPath)
	defer os.Remove(_logPath)

	logger.Println("dwd")

	logFile.Close()

	if _, err = os.Stat(_logPath); err != nil {
		if os.IsNotExist(err) {
			t.Error("File not exist: " + err.Error())
		} else {
			t.Error("File not allowed: " + err.Error())
		}
	} else {
		var file *os.File
		file, err = os.OpenFile(_logPath, os.O_RDWR, 0666)
		defer file.Close()

		if err != nil {
			t.Error("File not allowed: " + err.Error())
		}

		buffer := bytes.Buffer{}
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			buffer.WriteString(scanner.Text())
		}

		logSplited := strings.Split(buffer.String(), " ")
		log := logSplited[len(logSplited)-1]
		if log != "dwd" {
			t.Fatal("Log must have [dwd] but have [" + buffer.String() + "]")
		}
	}
}

func TestDB(t *testing.T) {
	var err error
	imgName := "test"
	imgExtension := "json"
	imgData := "Jojo bizare adventure"

	logger, logFile := CreateLog(_logPath)
	defer os.Remove(_logPath)
	defer logFile.Close()

	db := CreateDB(logger, &Configuration{DataBasePath: _dbPath})
	defer os.Remove(_dbPath)

	if err = db.CreateTable(); err != nil {
		t.Error("Table cant create: " + err.Error())
	}

	var id int64
	if id, err = db.SaveImage(imgName, imgExtension, imgData); err != nil {
		t.Error("Image not saved: " + err.Error())
	}

	var image *dto.Image
	if image, err = db.RestoreImage(id); err != nil {
		t.Error("Image did not restore: " + err.Error())
	}

	if image.Name != imgName {
		t.Fatal("Restore Name:[" + image.Name + "] not equal seved [" + imgName + "}")
	}

	if image.Extension != imgExtension {
		t.Fatal("Restore Extension:[" + image.Extension + "] not equal seved [" + imgExtension + "}")
	}

	if image.Data != imgData {
		t.Fatal("Restore Data:[" + image.Data + "] not equal seved [" + imgData + "}")
	}

	if err = db.Close(); err != nil {
		t.Error("DB not closed: " + err.Error())
	}
}

func TestValidator(t *testing.T) {
	logger, logFile := CreateLog(_logPath)
	defer os.Remove(_logPath)
	defer logFile.Close()

	validator := CreateValidator(logger, &Configuration{ScaledImageRestoreExtension: _extensionsCollection, FileSaveExtensionList: _extensionsCollection})

	for _, ext := range strings.Split(_extensionsCollection, ",") {
		if !validator.ValidateSavedFileExtension(ext) {
			t.Fatal(ext + " must be valid but not")
		}
	}

	if validator.ValidateSavedFileExtension("but") {
		t.Fatal("bat is not valid but true")
	}

	for _, ext := range strings.Split(_extensionsCollection, ",") {
		if !validator.ValidateScaledFileExtension(ext) {
			t.Fatal(ext + " must be valid but not")
		}
	}

	if validator.ValidateScaledFileExtension("but") {
		t.Fatal("bat is not valid but true")
	}

}
