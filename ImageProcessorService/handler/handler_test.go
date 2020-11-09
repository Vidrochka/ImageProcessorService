package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler/dto"
	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler/utils"
)

const _cfgPath string = "./test_config.toml"
const _logPath string = "./test.log"
const _dbPath string = "./test_db.db"
const _extensionsCollection string = "jpg,png,jpeg"

func CreateUtils() (*log.Logger, *os.File, *utils.DataBase, *utils.Configuration, *utils.Validator, *utils.FileSaver) {
	logger, logFile := utils.CreateLog(_logPath)

	config := &utils.Configuration{
		FileSaveExtensionList:       _extensionsCollection,
		ScaledImageRestoreExtension: _extensionsCollection,
		DataBasePath:                _dbPath,
		FileSavePath:                "./FileStorage/",
		PreviewFileFolder:           "Preview",
		ServedURL:                   "/",
		ScaledImageH:                100,
		ScaledImageW:                100,
	}

	db := utils.CreateDB(logger, config)
	if err := db.CreateTable(); err != nil {
		panic(err)
	}

	validator := utils.CreateValidator(logger, config)

	saver := utils.CreateFileSaver(logger, config)

	return logger, logFile, db, config, validator, saver
}

func TestBase64Handler(t *testing.T) {
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	buff := "{\"images\":[{\"name\": \"ff14\", \"extension\": \"png\", \"Data\": \"am9qbw==\"},{\"name\": \"ff14\", \"extension\": \"exe\", \"Data\": \"am9qbw==\"}]}"
	r := httptest.NewRequest("POST", "/", strings.NewReader(buff))
	w := httptest.NewRecorder()

	handler := CreateBase64(log, db, cfg, valid, saver)
	handler.Work(w, r)

	defer os.RemoveAll(cfg.FileSavePath)

	response := dto.Response{}
	if responseByte, err := ioutil.ReadAll(w.Body); err != nil {
		t.Error(err)
	} else {
		t.Log(string(responseByte))

		if err := json.Unmarshal(responseByte, &response); err != nil {
			t.Error(err)
		}
	}

	imageCollection := dto.ImageCollectionResponse{File: []dto.SaveImageResponseFile{}}

	messageBuff := bytes.NewBufferString(response.Message)
	if err := json.Unmarshal(messageBuff.Bytes(), &imageCollection); err != nil {
		t.Error(err)
	}

	if imageCollection.File[0].Name != "ff14" {
		t.Fatalf("File name not Equal. ff14 != %s", imageCollection.File[0].Name)
	}

	if imageCollection.File[0].Extension != "png" {
		t.Fatalf("File extension not Equal. png != %s", imageCollection.File[0].Extension)
	}

	if imageCollection.File[0].Status != 1 {
		t.Fatal("File extension not Equal. 1 != " + fmt.Sprint(imageCollection.File[0].Status))
	}

	t.Log("Valid Extension pass")

	if imageCollection.File[1].Name != "ff14" {
		t.Fatalf("File name not Equal. ff14 != %s", imageCollection.File[0].Name)
	}

	if imageCollection.File[1].Extension != "exe" {
		t.Fatalf("File extension not Equal. png != %s", imageCollection.File[0].Extension)
	}

	if imageCollection.File[1].Status != 0 {
		t.Fatal("File extension not Equal. 1 != " + fmt.Sprint(imageCollection.File[0].Status))
	}

	t.Log("Unvalid Extension pass")
}

func TestMultipatrFormDataHandler(t *testing.T) {
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	buff := "--X-INSOMNIA-BOUNDARY\r\n" +
		"Content-Disposition: form-data; name=\"image\"; filename=\"Новый текстовый документ (2).jpg\"\r\n" +
		"Content-Type: image/jpeg\r\n\r\n" +
		"gg\r\n" +
		"--X-INSOMNIA-BOUNDARY\r\n" +
		"Content-Disposition: form-data; name=\"image\"; filename=\"Новый текстовый документ.raw\"\r\n" +
		"Content-Type: image/raw\r\n\r\n" +
		"hh\r\n" +
		"--X-INSOMNIA-BOUNDARY--"
	r := httptest.NewRequest("POST", "/", strings.NewReader(buff))
	r.Header.Add("Content-Type", "multipart/form-data; boundary=X-INSOMNIA-BOUNDARY")
	r.ContentLength = 358
	w := httptest.NewRecorder()

	handler := CreateMultipartFormDataHandler(log, db, cfg, valid, saver)
	handler.Work(w, r)

	response := dto.Response{}
	if responseByte, err := ioutil.ReadAll(w.Body); err != nil {
		t.Error(err)
	} else {
		t.Log(string(responseByte))

		if err := json.Unmarshal(responseByte, &response); err != nil {
			t.Error(err)
		}
	}

	imageCollection := dto.ImageCollectionResponse{File: []dto.SaveImageResponseFile{}}
	if err := json.Unmarshal([]byte(response.Message), &imageCollection); err != nil {
		t.Error(err)
	}

	if imageCollection.File[0].Name != "Новый текстовый документ (2)" {
		t.Fatalf("File name not Equal. ff14 != %s", imageCollection.File[0].Name)
	}

	if imageCollection.File[0].Extension != "jpg" {
		t.Fatalf("File extension not Equal. png != %s", imageCollection.File[0].Extension)
	}

	if imageCollection.File[0].Status != 1 {
		t.Fatal("File extension not Equal. 1 != " + fmt.Sprint(imageCollection.File[0].Status))
	}

	t.Log("Valid Extension pass")

	if imageCollection.File[1].Name != "Новый текстовый документ" {
		t.Fatalf("File name not Equal. ff14 != %s", imageCollection.File[0].Name)
	}

	if imageCollection.File[1].Extension != "raw" {
		t.Fatalf("File extension not Equal. raw != %s", imageCollection.File[0].Extension)
	}

	if imageCollection.File[1].Status != 0 {
		t.Fatal("File extension not Equal. 1 != " + fmt.Sprint(imageCollection.File[0].Status))
	}

	t.Log("Unvalid Extension pass")

	os.RemoveAll(cfg.FileSavePath)
}

func TestRestoreImageHandler(t *testing.T) {
	fileName := "Fufix"
	fileExtension := "sos"
	fileData := "jojo-dwd-ara-ara"

	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	hash := saver.SaveFile(fileName, fileExtension, fileData)

	var buff string
	if id, err := db.SaveImage(fileName, fileExtension, hash); err != nil {
		t.Error(err)
	} else {
		buff = "{\"id\": " + fmt.Sprint(id) + "}"
	}

	r := httptest.NewRequest("POST", "/", strings.NewReader(buff))
	w := httptest.NewRecorder()

	handler := CreateRestore(log, db, valid, saver)
	handler.Work(w, r)

	response := dto.Response{}
	if responseByte, err := ioutil.ReadAll(w.Body); err != nil {
		t.Error(err)
	} else {
		t.Log(string(responseByte))

		if err = json.Unmarshal(responseByte, &response); err != nil {
			t.Error(err)
		}
	}

	image := dto.Image{}
	if data, err := base64.StdEncoding.DecodeString(response.Message); err != nil {
		t.Error(err)
	} else {
		if err = json.Unmarshal(data, &image); err != nil {
			t.Error(err)
		}
	}

	if image.Name != fileName {
		t.Fatalf("File name not Equal. %s != %s", fileName, image.Name)
	}

	if image.Extension != fileExtension {
		t.Fatalf("File extension not Equal. %s != %s", fileExtension, image.Extension)
	}

	if image.Data != fileData {
		t.Fatalf("File data not Equal. %s != %s", fileData, image.Data)
	}

	t.Log("Valid Extension pass")

	os.RemoveAll(cfg.FileSavePath)
}

func TestUrlLoadHandler(t *testing.T) {
	fileName := "peacock_PNG42"
	fileExtension := "png"

	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	buff := "{\"url\": \"http://pngimg.com/uploads/peacock/peacock_PNG42.png\"}"
	r := httptest.NewRequest("POST", "/", strings.NewReader(buff))
	w := httptest.NewRecorder()

	handler := CreateURLLoader(log, db, cfg, valid, saver)
	handler.Work(w, r)

	defer os.RemoveAll(cfg.FileSavePath)

	response := dto.Response{}
	if responseByte, err := ioutil.ReadAll(w.Body); err != nil {
		t.Error(err)
	} else {
		t.Log(string(responseByte))

		if err = json.Unmarshal(responseByte, &response); err != nil {
			t.Error(err)
		}
	}

	image := dto.SaveImageResponseFile{}
	if err := json.Unmarshal([]byte(response.Message), &image); err != nil {
		t.Error(err)
	}

	if image.Name != fileName {
		t.Fatalf("File name not Equal. %s != %s", fileName, image.Name)
	}

	if image.Extension != fileExtension {
		t.Fatalf("File extension not Equal. %s != %s", fileExtension, image.Extension)
	}

	if image.Status != 1 {
		t.Fatalf("File status not Equal. %d != %d", 1, image.Status)
	}

	t.Log("Valid Extension pass")

	var data []byte
	if imageFromDB, err := db.RestoreImage(image.ID); err != nil {
		t.Error(err)
	} else {
		if file, err := os.Open(cfg.FileSavePath + imageFromDB.Data + "/peacock_PNG42.png"); err != nil {
			t.Error(err)
		} else {
			defer file.Close()

			if data, err = ioutil.ReadAll(file); err != nil {
				t.Error(err)
			}
		}
	}

	if len(data) == 0 {
		t.Fatal("There is nothng data")
	}
}

func TestRestorePreviewImageHandler(t *testing.T) {
	urlPath := "https://steemitimages.com/DQmfWCZqLByZrs9oGSmYtH97iRwuwsuDsoZvK9EEpfXwxV7/"
	fileName := "pexels-photo-417074"
	fileExtension := "jpeg"

	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	buff := "{\"url\": \"" + urlPath + fileName + "." + fileExtension + "\"}"
	r := httptest.NewRequest("POST", "/", strings.NewReader(buff))
	w := httptest.NewRecorder()

	saveImageHandler := CreateURLLoader(log, db, cfg, valid, saver)
	saveImageHandler.Work(w, r)

	saveImageResponse := dto.Response{}
	if saveImageResponseByte, err := ioutil.ReadAll(w.Body); err != nil {
		t.Error(err)
	} else {
		t.Log(string(saveImageResponseByte))

		if err = json.Unmarshal(saveImageResponseByte, &saveImageResponse); err != nil {
			t.Error(err)
		}
	}

	saveImage := dto.SaveImageResponseFile{}
	if err := json.Unmarshal([]byte(saveImageResponse.Message), &saveImage); err != nil {
		t.Error(err)
	}

	buff = "{\"id\": " + fmt.Sprint(saveImage.ID) + "}"
	r = httptest.NewRequest("POST", "/", strings.NewReader(buff))
	w = httptest.NewRecorder()

	handler := CreatePrevievImageHandler(log, db, cfg, valid, saver)
	handler.Work(w, r)

	defer os.RemoveAll(cfg.FileSavePath)

	response := dto.Response{}
	if responseByte, err := ioutil.ReadAll(w.Body); err != nil {
		t.Error(err)
	} else {
		t.Log(string(responseByte))

		if err = json.Unmarshal(responseByte, &response); err != nil {
			t.Error(err)
		}

		t.Log(response.Message)
	}

	image := dto.Image{}
	if data, err := base64.StdEncoding.DecodeString(response.Message); err != nil {
		t.Error(err)
	} else {
		if err = json.Unmarshal(data, &image); err != nil {
			t.Error()
		}
	}

	if image.Name != fileName {
		t.Fatalf("File name not Equal. %s != %s", fileName, image.Name)
	}

	if image.Extension != fileExtension {
		t.Fatalf("File extension not Equal. %s != %s", fileExtension, image.Extension)
	}

	var genFileData []byte
	if imageFromDB, err := db.RestoreImage(saveImage.ID); err != nil {
		t.Error(err)
	} else {
		if file, err := os.Open(cfg.FileSavePath + imageFromDB.Data + "/" + cfg.PreviewFileFolder + "/" + fileName + "." + fileExtension + ""); err != nil {
			t.Error(err)
		} else {
			if genFileData, err = ioutil.ReadAll(file); err != nil {
				t.Error(err)
			}

			file.Close()
		}
	}

	var requiredData []byte
	if file, err := os.Open("./test_data/test_previev_image.jpeg"); err != nil {
		t.Error(err)
	} else {
		if requiredData, err = ioutil.ReadAll(file); err != nil {
			t.Error(err)
		}

		file.Close()
	}

	if string(genFileData) != string(requiredData) {
		t.Fatal("File data not Equal")
	}

	t.Log("Valid Extension pass")
}

func TestSelectorBadUrl(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("POST", "/dw", nil)
	r.Method = http.MethodPost
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Req-type", "BASE64")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)))
		t.Fatal("Must be bad request")
	}
}

func TestSelectorEmptyContentTypeRequest(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("POST", "/", nil)
	r.Method = http.MethodPost
	r.Header.Add("Req-type", "BASE64")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)))
		t.Fatal("Must be bad request")
	}
}

func TestSelectorBadMethidRequest(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Req-type", "BASE64")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)))
		t.Fatal("Must be bad request")
	}
}

func TestSelectorBase64Handler(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("POST", "/", nil)
	r.Method = http.MethodPost
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Req-type", "BASE64")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateBase64(log, db, cfg, valid, saver)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateBase64(log, db, cfg, valid, saver)))
		t.Fatal("Must be base64 request")
	}
}

func TestSelectorUrlHandler(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("POST", "/", nil)
	r.Method = http.MethodPost
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Req-type", "URL-LOAD")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateURLLoader(log, db, cfg, valid, saver)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateURLLoader(log, db, cfg, valid, saver)))
		t.Fatal("Must be url request")
	}
}

func TestSelectorRestoreHandler(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("POST", "/", nil)
	r.Method = http.MethodPost
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Req-type", "RESTORE")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateRestore(log, db, valid, saver)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateRestore(log, db, valid, saver)))
		t.Fatal("Must be restore request")
	}
}

func TestSelectorPreviewHandler(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("POST", "/", nil)
	r.Method = http.MethodPost
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Req-type", "RESTORE-PREVIEW")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreatePrevievImageHandler(log, db, cfg, valid, saver)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreatePrevievImageHandler(log, db, cfg, valid, saver)))
		t.Fatal("Must be multipart request")
	}
}

func TestSelectorMultipartHandler(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("POST", "/", nil)
	r.Method = http.MethodPost
	r.Header.Add("Content-type", "multipart/form-data")
	r.Header.Add("Req-type", "blablabla")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateMultipartFormDataHandler(log, db, cfg, valid, saver)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateMultipartFormDataHandler(log, db, cfg, valid, saver)))
		t.Fatal("Must be multipart request")
	}
}

func TestSelectorBadReqType(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Req-type", "BASE645653")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)))
		t.Fatal("Must be bad request")
	}
}

func TestSelectorBadContentTypeRequest(t *testing.T) {
	var err error
	log, logFile, db, cfg, valid, saver := CreateUtils()

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	defer os.Remove(_dbPath)
	defer db.Close()

	selector := CreateSelector(log, db, cfg, valid, saver)

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-type", "application")
	r.Header.Add("Req-type", "BASE64")
	w := httptest.NewRecorder()

	handler := selector.Select(w, r)
	if reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)) {
		handler.Work(w, r)
		var response []byte
		if response, err = ioutil.ReadAll(w.Body); err != nil {
			t.Error(err)
		}
		t.Log(string(response))
		t.Log(reflect.TypeOf(handler) != reflect.TypeOf(CreateBadRequestHandler("", "", 0, log)))
		t.Fatal("Must be bad request")
	}
}
