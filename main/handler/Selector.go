package handler

import (
	"ImageProcessorService/main/handler/utils"
	"fmt"
	"log"
	"net/http"
)

//Selector - select handler
type Selector struct {
	logger *log.Logger
	db     *utils.DataBase
	config *utils.Configuration
}

//CreateSelector - create selector
func CreateSelector(logger *log.Logger, db *utils.DataBase, config *utils.Configuration) *Selector {
	instanse := &Selector{logger: logger, db: db, config: config}

	logger.Println("Selector created")

	return instanse
}

//ServeHTTP - select handler
func (selector *Selector) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	selector.logger.Println("Has request")

	if req.RequestURI != selector.config.ServedURL {
		selector.logger.Println("Not supported url: " + req.RequestURI + "| need: " + selector.config.ServedURL)
		resp.WriteHeader(405)
		fmt.Fprint(resp, "Not supported url: "+req.RequestURI+"| need: "+selector.config.ServedURL)
		return
	}

	if req.Method != http.MethodPost {
		selector.logger.Println("Not Post request")
		resp.WriteHeader(405)
		fmt.Fprint(resp, "Served only post request")
		return
	}

	contentType := req.Header.Get("Content-type")
	selector.logger.Printf("Content-type = %s", contentType)

	if contentType == "" {
		selector.logger.Println("Request without Content-type")
		fmt.Fprint(resp, "You need set Content-type")
		return
	}

	var handler Handler = nil

	switch contentType {
	case "application/json":
		switch req.Header.Get("Req-type") {
		case "BASE64":
			handler = CreateBase64(selector.logger, selector.db)
			break
		case "URL-LOAD":
			handler = CreateURLLoader(selector.logger, selector.db, []string{"jpg"})
			break
		case "RESTORE":
			handler = CreateRestore(selector.logger, selector.db)
			break
		case "RESTORE-PREVIEW":
			handler = CreateSquareImageHandler(selector.logger, selector.db, []string{"jpg"})
			break
		default:
			fmt.Fprint(resp, "Served only request with header Req-type = BASE64/URL-LOAD/RESTORE/RESTORE-PREVIEW")
			return
		}
		break
	case "multipart/form-data":
		break
	default:
		fmt.Fprint(resp, "Served only request with application/json or multipart/form-data Content-type")
		return
	}

	if handler == nil {
		panic("handler is not defined!!!")
	}

	handler.Work(resp, req)
}
