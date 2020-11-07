package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler"
	"github.com/anthill-com/ImageProcessorService/ImageProcessorService/handler/utils"
)

//Server - server wrapper
type Server struct {
	httpServer *http.Server
	logger     *log.Logger
	config     *utils.Configuration
	selector   *handler.Selector
}

//CreateServer - create server
func CreateServer(logger *log.Logger, config *utils.Configuration, selector *handler.Selector) *Server {
	server := Server{logger: logger, config: config, selector: selector}

	logger.Println("Server created")

	return &server
}

//Run - run server
func (server *Server) Run() error {
	server.logger.Println("Listen port: " + server.config.Port)
	server.logger.Println("Read timeout: " + fmt.Sprint(server.config.ReadTimeout) + "sec")
	server.logger.Println("Write timeout: " + fmt.Sprint(server.config.WriteTimeout) + "sec")

	server.httpServer = &http.Server{
		Addr:           ":" + server.config.Port,
		MaxHeaderBytes: 1 << 28,
		ReadTimeout:    server.config.ReadTimeout * time.Second,
		WriteTimeout:   server.config.WriteTimeout * time.Second,
	}

	server.httpServer.Handler = server.selector

	return server.httpServer.ListenAndServe()
}

//Stop - stop server
func (server *Server) Stop(ctx context.Context) error {
	return server.httpServer.Shutdown(ctx)
}
