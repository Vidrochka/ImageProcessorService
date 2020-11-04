package handler

import (
	"ImageProcessorService/main/handler/dto"
	"ImageProcessorService/main/handler/utils"
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

	"github.com/nfnt/resize"
)

//SquareImageHandler - make square image
type SquareImageHandler struct {
	logger            *log.Logger
	db                *utils.DataBase
	supportExtensions []string
}

//CreateSquareImageHandler - create url loader handler
func CreateSquareImageHandler(logger *log.Logger, db *utils.DataBase, supportedExtension []string) Handler {
	var instanse Handler = &SquareImageHandler{logger: logger, db: db, supportExtensions: supportedExtension}

	logger.Println("Square image handler created")

	return instanse
}

//Work - work with square image request
func (handler *SquareImageHandler) Work(resp http.ResponseWriter, req *http.Request) {
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

	var restoreImage *dto.Image
	restoreImage, err = handler.db.RestoreImage(request.ID)

	if err != nil {
		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant find image", ResCode: 1}.ToJSON())
		return
	}

	if !handler.IsSupportImage(restoreImage.Extension) {
		handler.logger.Print("Not supported extension - " + restoreImage.Extension + " | valid: " + strings.Join(handler.supportExtensions, "/"))
		resp.WriteHeader(415)
		fmt.Fprintf(resp, dto.Response{Message: "Not supported extension - " + restoreImage.Extension + " | valid: " + strings.Join(handler.supportExtensions, "/"), ResCode: 2}.ToJSON())
		return
	}

	var imageData image.Image
	switch restoreImage.Extension {
	case "jpeg":
		imageData, err = jpeg.Decode(strings.NewReader(restoreImage.Data))
		break
	case "jpg":
		imageData, err = jpeg.Decode(strings.NewReader(restoreImage.Data))
		break
	case "png":
		imageData, err = png.Decode(strings.NewReader(restoreImage.Data))
		break
	default:
		panic("There is no handler for that extansion, fix that")
	}

	resizedImage := resize.Resize(100, 100, imageData, resize.Lanczos3)

	buf := bytes.NewBufferString("")

	err = jpeg.Encode(buf, resizedImage, nil)

	if err != nil {

		handler.logger.Println(err)
		resp.WriteHeader(418)
		fmt.Fprintf(resp, dto.Response{Message: "Sorry, we cant make image smaller", ResCode: 1}.ToJSON())
		return
	}

	restoreImage.Data = buf.String()

	response := restoreImage.ToJSON()
	handler.logger.Print(response)
	resp.WriteHeader(200)
	fmt.Fprintf(resp, dto.Response{Message: response, ResCode: 0}.ToJSON())
}

//IsSupportImage - check if extension is supported
func (handler *SquareImageHandler) IsSupportImage(extension string) bool {
	for _, ext := range handler.supportExtensions {
		if ext == extension {
			return true
		}
	}

	return false
}
