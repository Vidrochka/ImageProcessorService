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
	logger    *log.Logger
	db        *utils.DataBase
	validator *utils.Validator
}

//CreateRestore - create restore hendler
func CreateRestore(logger *log.Logger, db *utils.DataBase, validator *utils.Validator) Handler {
	var instanse Handler = &RestoreImageHandler{logger: logger, db: db, validator: validator}

	logger.Println("Restore handler created")

	return instanse
}

//Work - implement Handler interfase
func (handler *RestoreImageHandler) Work(resp http.ResponseWriter, req *http.Request) {
	var err error

	var data []byte
	if data, err = ioutil.ReadAll(req.Body); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request", ResCode: 1}.ToJSON())
		return
	}

	var request dto.ImageIdRequest

	if err = json.Unmarshal(data, &request); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Invalid Json format \"%s" + err.Error() + "\"", ResCode: 2}.ToJSON())
		return
	}

	var image *dto.Image
	if image, err = handler.db.RestoreImage(request.ID); err != nil {
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
