package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/anthill-com/ImageProcessorService/main/handler/dto"
	"github.com/anthill-com/ImageProcessorService/main/handler/utils"
)

//URLLoadHandler - load image by url
type URLLoadHandler struct {
	logger    *log.Logger
	db        *utils.DataBase
	config    *utils.Configuration
	validator *utils.Validator
	fileSaver *utils.FileSaver
}

//CreateURLLoader - create url loader handler
func CreateURLLoader(logger *log.Logger, db *utils.DataBase, config *utils.Configuration, validator *utils.Validator, fileSaver *utils.FileSaver) Handler {
	var instanse Handler = &URLLoadHandler{logger: logger, db: db, config: config, validator: validator, fileSaver: fileSaver}

	logger.Println("Url load handler created")

	return instanse
}

//Work - implement Handler interfase
func (handler *URLLoadHandler) Work(resp http.ResponseWriter, req *http.Request) {
	var err error

	resp.Header().Set("Content-Type", "application/json; charset=utf-8")

	var data []byte
	if data, err = ioutil.ReadAll(req.Body); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request", ResCode: 1}.ToJSON())
		return
	}

	var request dto.URLLoadRequest
	if err = json.Unmarshal(data, &request); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Invalid Json format \"%s" + err.Error() + "\"", ResCode: 2}.ToJSON())
		return
	}

	pathURLSplited := strings.Split(request.URL, "/")
	dotIndex := strings.LastIndex(pathURLSplited[len(pathURLSplited)-1], ".")

	if dotIndex == -1 {
		handler.logger.Printf("Url not contain extension - %s", request.URL)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Url not contain extension - %s" + request.URL, ResCode: 2}.ToJSON())
		return
	}

	name := pathURLSplited[len(pathURLSplited)-1][:dotIndex]
	extension := pathURLSplited[len(pathURLSplited)-1][dotIndex+1:]

	handler.logger.Printf("File name: %s | File extension: %s", name, extension)

	if !handler.validator.ValidateSavedFileExtension(extension) {
		handler.logger.Print("Not supported extension - " + extension + " | valid: " + handler.config.ScaledImageRestoreExtension)
		resp.WriteHeader(415)
		fmt.Fprintf(resp, dto.Response{Message: "Not supported extension - " + extension + " | valid: " + handler.config.ScaledImageRestoreExtension, ResCode: 2}.ToJSON())
		return
	}

	var imageResponse *http.Response
	if imageResponse, err = http.Get(request.URL); err != nil || imageResponse.StatusCode != http.StatusOK {
		handler.logger.Println(err)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request by url - " + request.URL, ResCode: 1}.ToJSON())
		return
	}

	if data, err = ioutil.ReadAll(imageResponse.Body); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request by url - " + request.URL, ResCode: 1}.ToJSON())
		return
	}

	hash := handler.fileSaver.SaveFile(name, extension, string(data))
	var id int64
	if id, err = handler.db.SaveImage(name, extension, hash); err != nil {
		handler.logger.Printf("Image not saved - %s, error - %s", name, err.Error())
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant save image name = " + name + " extension = " + extension, ResCode: 1}.ToJSON())
		return
	}

	handler.logger.Printf("Image saved - %s, id - %d", name, id)
	resp.WriteHeader(200)
	fmt.Fprintf(resp, dto.Response{Message: dto.SaveImageResponseFile{ID: id, Name: name, Extension: extension, Status: 1}.ToJSON(), ResCode: 0}.ToJSON())

}
