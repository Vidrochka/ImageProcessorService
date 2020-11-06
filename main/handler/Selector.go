package handler

import (
	"ImageProcessorService/main/handler/dto"
	"ImageProcessorService/main/handler/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//Selector - select handler
type Selector struct {
	logger    *log.Logger
	db        *utils.DataBase
	config    *utils.Configuration
	validator *utils.Validator
}

//CreateSelector - create selector
func CreateSelector(logger *log.Logger, db *utils.DataBase, config *utils.Configuration, validator *utils.Validator) *Selector {
	instanse := &Selector{logger: logger, db: db, config: config, validator: validator}

	logger.Println("Selector created")

	return instanse
}

//ServeHTTP - select handler
func (selector *Selector) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	selector.logger.Println("Has request")

	if req.RequestURI != selector.config.ServedURL {
		selector.logger.Println("Not supported url: " + req.RequestURI + "| need: " + selector.config.ServedURL)
		resp.WriteHeader(405)
		fmt.Fprint(resp, dto.Response{Message: "Not supported url: " + req.RequestURI + "| need: " + selector.config.ServedURL, ResCode: 0}.ToJSON())
		return
	}

	if req.Method != http.MethodPost {
		selector.logger.Println("Not Post request")
		resp.WriteHeader(405)
		fmt.Fprint(resp, dto.Response{Message: "Served only post request", ResCode: 0}.ToJSON())
		return
	}

	contentType := req.Header.Get("Content-type")
	selector.logger.Printf("Content-type = %s", contentType)

	if contentType == "" {
		selector.logger.Println("Request without Content-type")
		fmt.Fprint(resp, dto.Response{Message: "You need set Content-type", ResCode: 0}.ToJSON())
		return
	}

	var handler Handler = nil

	if contentType == "application/json" {
		switch req.Header.Get("Req-type") {
		case "BASE64":
			handler = CreateBase64(selector.logger, selector.db, selector.config, selector.validator)
			break
		case "URL-LOAD":
			handler = CreateURLLoader(selector.logger, selector.db, selector.config, selector.validator)
			break
		case "RESTORE":
			handler = CreateRestore(selector.logger, selector.db, selector.validator)
			break
		case "RESTORE-PREVIEW":
			handler = CreateScaleImageHandler(selector.logger, selector.db, selector.config, selector.validator)
			break
		default:
			fmt.Fprint(resp, dto.Response{Message: "Served only request with header Req-type = BASE64/URL-LOAD/RESTORE/RESTORE-PREVIEW", ResCode: 0}.ToJSON())
			return
		}
	} else {
		if strings.Contains(contentType, "multipart/form-data") {
			handler = CreateMultipartFormDataHandler(selector.logger, selector.db, selector.config, selector.validator)
		} else {
			fmt.Fprint(resp, dto.Response{Message: "Served only request with application/json or multipart/form-data Content-type", ResCode: 2}.ToJSON())
			return
		}
	}

	if handler == nil {
		panic("handler is not defined!!!")
	}

	handler.Work(resp, req)
}
