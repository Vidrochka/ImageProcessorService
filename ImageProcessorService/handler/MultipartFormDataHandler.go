package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler/dto"
	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler/utils"
)

//MultipartFormDataHandler - handler which parse multipart\form-data
type MultipartFormDataHandler struct {
	logger    *log.Logger
	db        *utils.DataBase
	validator *utils.Validator
	config    *utils.Configuration
	fileSaver *utils.FileSaver
}

//CreateMultipartFormDataHandler - create restore hendler
func CreateMultipartFormDataHandler(logger *log.Logger, db *utils.DataBase, config *utils.Configuration, validator *utils.Validator, fileSaver *utils.FileSaver) Handler {
	var instanse Handler = &MultipartFormDataHandler{logger: logger, db: db, config: config, validator: validator, fileSaver: fileSaver}

	logger.Println("Restore handler created")

	return instanse
}

//Work - work with square image request
func (handler *MultipartFormDataHandler) Work(resp http.ResponseWriter, req *http.Request) {
	var err error

	if err = req.ParseMultipartForm(0); err != nil {
		handler.logger.Println("We cant parse multipart: " + err.Error())
		resp.WriteHeader(400)
		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(resp, dto.Response{Message: "We cant parse multipart", ResCode: 2}.ToJSON())
		return
	}

	formdata := req.MultipartForm
	files := formdata.File["image"]

	resp.Header().Set("Content-Type", "application/json; charset=utf-8")

	var requestCollection []dto.SaveImageResponseFile

	for i, _ := range files {

		var file multipart.File
		if file, err = files[i].Open(); err != nil {
			handler.logger.Println(err)
			resp.WriteHeader(400)
			fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request", ResCode: 1}.ToJSON())
			return
		}
		defer file.Close()

		fileNameSplitter := strings.LastIndex(files[i].Filename, ".")

		fileName := files[i].Filename[:fileNameSplitter]
		fileExtension := files[i].Filename[fileNameSplitter+1:]

		if !handler.validator.ValidateScaledFileExtension(fileExtension) {
			handler.logger.Println("Not valid extension: " + fileExtension)
			requestCollection = append(requestCollection, dto.SaveImageResponseFile{ID: -1, Name: fileName, Extension: fileExtension, Status: 0, ResMessage: "Not valid extension: " + fileExtension + " | valid: " + handler.config.FileSaveExtensionList})
			continue
		}

		buffer := bytes.Buffer{}
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			buffer.WriteString(scanner.Text())
		}

		hash := handler.fileSaver.SaveFile(fileName, fileExtension, buffer.String())
		var id int64
		if id, err = handler.db.SaveImage(fileName, fileExtension, hash); err != nil {
			handler.logger.Printf("Image not saved - %s, error - %s", fileName, err.Error())
			continue
		} else {
			handler.logger.Printf("Image saved - %s", fileName)
		}

		requestCollection = append(requestCollection, dto.SaveImageResponseFile{ID: id, Name: fileName, Extension: fileExtension, Status: 1, ResMessage: ""})
	}

	var response string
	if requestCollection == nil {
		response = dto.Response{Message: "We cant write anything data", ResCode: 2}.ToJSON()
		handler.logger.Print(response)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, response)
		return
	}

	response = dto.ImageCollectionResponse{File: requestCollection}.ToJSON()
	handler.logger.Print(response)
	resp.WriteHeader(200)
	fmt.Fprintf(resp, dto.Response{Message: response, ResCode: 0}.ToJSON())
}
