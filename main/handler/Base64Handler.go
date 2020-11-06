package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"ImageProcessorService/main/handler/dto"
	"ImageProcessorService/main/handler/utils"
)

//Base64Handler - work with base64 incode image
type Base64Handler struct {
	logger    *log.Logger
	db        *utils.DataBase
	config    *utils.Configuration
	validator *utils.Validator
}

//CreateBase64 - create base64 request handler
func CreateBase64(logger *log.Logger, db *utils.DataBase, config *utils.Configuration, validator *utils.Validator) Handler {
	var instanse Handler = &Base64Handler{logger: logger, db: db, config: config, validator: validator}

	logger.Println("Base64 handler created")

	return instanse
}

//Work - implement Handler interfase
func (handler *Base64Handler) Work(resp http.ResponseWriter, req *http.Request) {
	var err error

	var data []byte
	if data, err = ioutil.ReadAll(req.Body); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request", ResCode: 1}.ToJSON())
		return
	}

	var request dto.Base64Request

	if err = json.Unmarshal(data, &request); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Invalid Json format \"%s" + err.Error() + "\"", ResCode: 2}.ToJSON())
		return
	}

	var imgDecodeStr []byte
	var requestCollection []dto.SaveImageResponseFile

	for _, file := range request.Images {
		var id int64

		if !handler.validator.ValidateScaledFileExtension(file.Extension) {
			handler.logger.Println("Not valid extension: " + file.Extension)
			requestCollection = append(requestCollection, dto.SaveImageResponseFile{ID: -1, Name: file.Name, Extension: file.Extension, Status: 0, ResMessage: "Not valid extension: " + file.Extension + " | valid: " + handler.config.FileSaveExtensionList})
			continue
		}

		if imgDecodeStr, err = base64.StdEncoding.DecodeString(file.Data); err != nil {
			handler.logger.Println(err)
			requestCollection = append(requestCollection, dto.SaveImageResponseFile{ID: -1, Name: file.Name, Extension: file.Extension, Status: 0, ResMessage: fmt.Sprintf("Invalid base64 format \"%s", err.Error())})
			continue
		}

		if id, err = handler.db.SaveImage(file.Name, file.Extension, string(imgDecodeStr)); err != nil {
			handler.logger.Printf("Image not saved - %s, error - %s", file.Name, err.Error())
			continue
		} else {
			handler.logger.Printf("Image saved - %s", file.Name)
		}

		requestCollection = append(requestCollection, dto.SaveImageResponseFile{ID: id, Name: file.Name, Extension: file.Extension, Status: 1, ResMessage: ""})
	}

	var response string
	if requestCollection == nil {
		response = dto.Response{Message: "We cant write anything data", ResCode: 2}.ToJSON()
		handler.logger.Print(response)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, response)
		return
	}

	response = dto.ImageCollection{File: requestCollection}.ToJSON()
	handler.logger.Print(response)
	resp.WriteHeader(200)
	fmt.Fprintf(resp, dto.Response{Message: response, ResCode: 0}.ToJSON())
}
