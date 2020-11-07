package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/anthill-com/ImageProcessorService/main/handler/dto"
	"github.com/anthill-com/ImageProcessorService/main/handler/utils"

	"github.com/nfnt/resize"
)

//PrevievImageHandler - make square image
type PrevievImageHandler struct {
	logger    *log.Logger
	db        *utils.DataBase
	config    *utils.Configuration
	validator *utils.Validator
	fileSaver *utils.FileSaver
}

//CreatePrevievImageHandler - create url loader handler
func CreatePrevievImageHandler(logger *log.Logger, db *utils.DataBase, config *utils.Configuration, validator *utils.Validator, fileSaver *utils.FileSaver) Handler {
	var instanse Handler = &PrevievImageHandler{logger: logger, db: db, config: config, validator: validator, fileSaver: fileSaver}

	logger.Println("Previev image handler created")

	return instanse
}

//Work - work with scale image request
func (handler *PrevievImageHandler) Work(resp http.ResponseWriter, req *http.Request) {
	var err error

	resp.Header().Set("Content-Type", "application/json; charset=utf-8")

	var data []byte
	if data, err = ioutil.ReadAll(req.Body); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant take request", ResCode: 1}.ToJSON())
		return
	}

	var request dto.ImageIDRequest
	if err = json.Unmarshal(data, &request); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(400)
		fmt.Fprintf(resp, dto.Response{Message: "Invalid Json format \"%s" + err.Error() + "\"", ResCode: 2}.ToJSON())
		return
	}

	var restoreImage *dto.Image
	if restoreImage, err = handler.db.RestoreImage(request.ID); err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant find image", ResCode: 1}.ToJSON())
		return
	}

	var hash string = restoreImage.Data
	handler.logger.Println("Hash: " + hash)

	restoreData := handler.fileSaver.RestoreFile(restoreImage.Name, restoreImage.Extension, restoreImage.Data)
	restoreImage.Data = restoreData

	if !handler.validator.ValidateScaledFileExtension(restoreImage.Extension) {
		handler.logger.Print("Not supported extension - " + restoreImage.Extension + " | valid: " + handler.config.ScaledImageRestoreExtension)
		resp.WriteHeader(415)
		fmt.Fprintf(resp, dto.Response{Message: "Not supported extension - " + restoreImage.Extension + " | valid: " + handler.config.ScaledImageRestoreExtension, ResCode: 2}.ToJSON())
		return
	}

	var imageData image.Image
	switch restoreImage.Extension {
	case "jpeg":
		if imageData, err = jpeg.Decode(strings.NewReader(restoreImage.Data)); err != nil {
			handler.logger.Println("Not valid jpeg data: " + err.Error())
			resp.WriteHeader(418)
			fmt.Fprintf(resp, dto.Response{Message: "Not valid jpeg data", ResCode: 2}.ToJSON())
			return
		}
		break
	case "jpg":
		if imageData, err = jpeg.Decode(strings.NewReader(restoreImage.Data)); err != nil {
			handler.logger.Println("Not valid jpg data: " + err.Error())
			resp.WriteHeader(418)
			fmt.Fprintf(resp, dto.Response{Message: "Not valid jpg data", ResCode: 2}.ToJSON())
			return
		}
		break
	case "png":
		if imageData, err = png.Decode(strings.NewReader(restoreImage.Data)); err != nil {
			handler.logger.Println("Not valid png data: " + err.Error())
			resp.WriteHeader(418)
			fmt.Fprintf(resp, dto.Response{Message: "Not valid png data", ResCode: 2}.ToJSON())
			return
		}
		break
	default:
		panic("There is no handler for that extansion, fix that")
	}

	handler.logger.Printf("Resize image to: w = %d, h = %d", handler.config.ScaledImageW, handler.config.ScaledImageH)
	resizedImage := resize.Resize(handler.config.ScaledImageW, handler.config.ScaledImageH, imageData, resize.Lanczos3)

	buf := bytes.NewBufferString("")

	if err = jpeg.Encode(buf, resizedImage, nil); err != nil {

		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant make image smaller", ResCode: 1}.ToJSON())
		return
	}

	restoreImage.Data = buf.String()

	handler.fileSaver.SavePreview(restoreImage.Name, restoreImage.Extension, restoreImage.Data, hash)

	response := restoreImage.ToJSON()
	handler.logger.Print(response)
	resp.WriteHeader(200)
	fmt.Fprintf(resp, dto.Response{Message: response, ResCode: 0}.ToJSON())
}
