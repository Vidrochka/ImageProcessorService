package handler

import (
	"ImageProcessorService/main/handler/dto"
	"ImageProcessorService/main/handler/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//RestoreImageHandler - restore image from db
type RestoreImageHandler struct {
	logger *log.Logger
	db     *utils.DataBase
}

//CreateRestore - create restore hendler
func CreateRestore(logger *log.Logger, db *utils.DataBase) Handler {
	var instanse Handler = &RestoreImageHandler{logger: logger, db: db}

	logger.Println("Restore handler created")

	return instanse
}

//Work - implement Handler interfase
func (handler *RestoreImageHandler) Work(resp http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)

	if err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request", ResCode: 1}.ToJSON())
		return
	}

	var request dto.ImageIdRequest

	err = json.Unmarshal(data, &request)

	if err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Invalid Json format \"%s" + err.Error() + "\"", ResCode: 2}.ToJSON())
		return
	}

	var image *dto.Image
	image, err = handler.db.RestoreImage(request.ID)

	if err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant restore image", ResCode: 1}.ToJSON())
		return
	}

	response := image.ToJSON()
	handler.logger.Print(response)
	resp.WriteHeader(200)
	fmt.Fprintf(resp, dto.Response{Message: response, ResCode: 0}.ToJSON())
}
