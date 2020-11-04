package handler

import (
	"ImageProcessorService/main/handler/dto"
	"ImageProcessorService/main/handler/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//URLLoadHandler - load image by url
type URLLoadHandler struct {
	logger            *log.Logger
	db                *utils.DataBase
	supportExtensions []string
}

//CreateURLLoader - create url loader handler
func CreateURLLoader(logger *log.Logger, db *utils.DataBase, supportedExtension []string) Handler {
	var instanse Handler = &URLLoadHandler{logger: logger, db: db, supportExtensions: supportedExtension}

	logger.Println("Url load handler created")

	return instanse
}

//Work - implement Handler interfase
func (handler *URLLoadHandler) Work(resp http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)

	if err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request", ResCode: 1}.ToJSON())
		return
	}

	var request dto.URLLoadRequest

	err = json.Unmarshal(data, &request)

	if err != nil {
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

	if !handler.IsSupportExtension(extension) {
		handler.logger.Print("Not supported extension - " + extension + " | valid: " + strings.Join(handler.supportExtensions, "/"))
		resp.WriteHeader(415)
		fmt.Fprintf(resp, dto.Response{Message: "Not supported extension - " + extension + " | valid: " + strings.Join(handler.supportExtensions, "/"), ResCode: 2}.ToJSON())
		return
	}

	var imageResponse *http.Response
	imageResponse, err = http.Get(request.URL)

	if err != nil || imageResponse.StatusCode != http.StatusOK {
		handler.logger.Println(err)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request by url - " + request.URL, ResCode: 1}.ToJSON())
		return
	}

	data, err = ioutil.ReadAll(imageResponse.Body)

	if err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request by url - " + request.URL, ResCode: 1}.ToJSON())
		return
	}

	var id int64

	id, err = handler.db.SaveImage(name, extension, string(data))

	if err != nil {
		handler.logger.Printf("Image not saved - %s, error - %s", name, err.Error())
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant save image name = " + name + " extension = " + extension, ResCode: 1}.ToJSON())
		return
	}

	handler.logger.Printf("Image saved - %s, id - %d", name, id)
	resp.WriteHeader(200)
	fmt.Fprintf(resp, dto.Response{Message: dto.SaveImageResponseFile{ID: id, Name: name}.ToJSON(), ResCode: 0}.ToJSON())

}

//IsSupportExtension - check if extension support
func (handler *URLLoadHandler) IsSupportExtension(extension string) bool {
	for _, ext := range handler.supportExtensions {
		if ext == extension {
			return true
		}
	}

	return false
}
