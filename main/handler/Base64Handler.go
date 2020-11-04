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
	logger *log.Logger
	db     *utils.DataBase
}

//CreateBase64 - create base64 request handler
func CreateBase64(logger *log.Logger, db *utils.DataBase) Handler {
	var instanse Handler = &Base64Handler{logger: logger, db: db}

	logger.Println("Base64 handler created")

	return instanse
}

//Work - implement Handler interfase
func (handler *Base64Handler) Work(resp http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)

	if err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request", ResCode: 1}.ToJSON())
		return
	}

	var request dto.Base64Request

	err = json.Unmarshal(data, &request)

	if err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Invalid Json format \"%s" + err.Error() + "\"", ResCode: 2}.ToJSON())
		return
	}

	var imgDecodeStr []byte
	var requestCollection []dto.SaveImageResponseFile

	for _, file := range request.Images {
		var id int64
		imgDecodeStr, err = base64.StdEncoding.DecodeString(file.Data)

		if err != nil {
			handler.logger.Println(err)
			resp.WriteHeader(400)
			fmt.Fprint(resp, dto.Response{Message: "Invalid base64 format \"%s" + err.Error() + "\"", ResCode: 2}.ToJSON())
			return
		}

		id, err = handler.db.SaveImage(file.Name, file.Extension, string(imgDecodeStr))

		if err != nil {
			handler.logger.Printf("Image not saved - %s, error - %s", file.Name, err.Error())
			continue
		} else {
			handler.logger.Printf("Image saved - %s", file.Name)
		}

		requestCollection = append(requestCollection, dto.SaveImageResponseFile{ID: id, Name: file.Name})
	}

	var response string
	if requestCollection == nil {
		response = dto.Response{Message: "We cant write anything data", ResCode: 2}.ToJSON()
		handler.logger.Print(response)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, response)
		return
	}
	response = dto.Base64Response{File: requestCollection}.ToJSON()
	handler.logger.Print(response)
	resp.WriteHeader(200)
	fmt.Fprintf(resp, dto.Response{Message: response, ResCode: 0}.ToJSON())
}
