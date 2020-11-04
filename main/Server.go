package main

import (
	"ImageProcessorService/main/handler"
	"ImageProcessorService/main/handler/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
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
