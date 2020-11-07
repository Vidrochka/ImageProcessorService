package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler/dto"
)

//RequestErrorHandler - bad request handler
type RequestErrorHandler struct {
	logger      *log.Logger
	logMessage  string
	respMessage string
	resCode     int
	httpCode    int
}

//CreateBadRequestHandler - create bad request handler
func CreateBadRequestHandler(logMessage, respMessage string, resCode int, logger *log.Logger) Handler {
	handler := RequestErrorHandler{logMessage: logMessage, respMessage: respMessage, resCode: resCode, httpCode: 400, logger: logger}

	logger.Println("BadRequestHandler created")

	return &handler
}

//CreateServerErrorHandler - create server error handler
func CreateServerErrorHandler(logMessage, respMessage string, resCode int, logger *log.Logger) Handler {
	handler := RequestErrorHandler{logMessage: logMessage, respMessage: respMessage, resCode: resCode, httpCode: 418, logger: logger}

	logger.Println("BadRequestHandler created")

	return &handler
}

//CreateErrorHandler - create error handler
func CreateErrorHandler(logMessage, respMessage string, resCode, httpCode int, logger *log.Logger) Handler {
	handler := RequestErrorHandler{logMessage: logMessage, respMessage: respMessage, resCode: resCode, httpCode: httpCode, logger: logger}

	logger.Println("BadRequestHandler created")

	return &handler
}

//Work - write bad request
func (handler *RequestErrorHandler) Work(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Println(handler.logMessage)

	resp.WriteHeader(handler.httpCode)

	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(resp, dto.Response{Message: handler.respMessage, ResCode: 0}.ToJSON())
}
