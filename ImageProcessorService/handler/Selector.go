package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler/utils"
)

//Selector - select handler
type Selector struct {
	logger    *log.Logger
	db        *utils.DataBase
	config    *utils.Configuration
	validator *utils.Validator
	fileSaver *utils.FileSaver
}

//CreateSelector - create selector
func CreateSelector(logger *log.Logger, db *utils.DataBase, config *utils.Configuration, validator *utils.Validator, fileSaver *utils.FileSaver) *Selector {
	instanse := &Selector{logger: logger, db: db, config: config, validator: validator, fileSaver: fileSaver}

	logger.Println("Selector created")

	return instanse
}

//ServeHTTP - select handler
func (selector *Selector) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	selector.logger.Println("Has request")

	handler := selector.Select(resp, req)

	if handler == nil {
		panic("handler is not defined!!!")
	}

	handler.Work(resp, req)
}

//Select - select handler
func (selector *Selector) Select(resp http.ResponseWriter, req *http.Request) Handler {
	if req.RequestURI != selector.config.ServedURL {
		return CreateBadRequestHandler(
			"Not supported url: "+req.RequestURI+"| need: "+selector.config.ServedURL,
			"Not supported url: "+req.RequestURI+"| need: "+selector.config.ServedURL,
			2,
			selector.logger)
	}

	if req.Method != http.MethodPost {
		return CreateErrorHandler(
			"Not Post request",
			"Served only post request",
			2,
			405,
			selector.logger)
	}

	contentType := req.Header.Get("Content-type")
	selector.logger.Printf("Content-type = %s", contentType)

	if contentType == "" {
		return CreateBadRequestHandler(
			"Request without Content-type",
			"You need set Content-type",
			2,
			selector.logger)
	}

	if contentType == "application/json" {
		switch req.Header.Get("Req-type") {
		case "BASE64":
			return CreateBase64(selector.logger, selector.db, selector.config, selector.validator, selector.fileSaver)
		case "URL-LOAD":
			return CreateURLLoader(selector.logger, selector.db, selector.config, selector.validator, selector.fileSaver)
		case "RESTORE":
			return CreateRestore(selector.logger, selector.db, selector.validator, selector.fileSaver)
		case "RESTORE-PREVIEW":
			return CreatePrevievImageHandler(selector.logger, selector.db, selector.config, selector.validator, selector.fileSaver)
		default:
			return CreateBadRequestHandler(
				"Incorrect Req-type"+req.Header.Get("Req-type"),
				"Served only request with header Req-type = BASE64/URL-LOAD/RESTORE/RESTORE-PREVIEW",
				2,
				selector.logger)
		}
	} else {
		if strings.Contains(contentType, "multipart/form-data") {
			return CreateMultipartFormDataHandler(selector.logger, selector.db, selector.config, selector.validator, selector.fileSaver)
		} else {
			return CreateBadRequestHandler(
				"Incorrect contentType: "+contentType,
				"Served only request with application/json or multipart/form-data Content-type",
				2,
				selector.logger)
		}
	}
}
