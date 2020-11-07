package handler

import (
	"net/http"
)

//Handler - handler base interface
type Handler interface {
	Work(http.ResponseWriter, *http.Request)
}
