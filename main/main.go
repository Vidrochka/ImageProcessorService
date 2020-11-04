package main

import (
	"ImageProcessorService/main/handler"
	"ImageProcessorService/main/handler/utils"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := utils.LoadConfiguration("config.toml")

	logger := utils.CreateLog(config.LogFilePath)

	dataBase := utils.CreateDB(logger)
	dataBase.CreateTable()

	selector := handler.CreateSelector(logger, dataBase, config)

	server := CreateServer(logger, config, selector)

	var err error

	go func() {
		if err = server.Run(); err != nil {
			logger.Println(err)
			return
		}
	}()

	logger.Println("Server starsed")

	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGTERM, syscall.SIGINT)

	<-quite

	if err = server.Stop(context.Background()); err != nil {
		logger.Println(err)
	}

	if err = dataBase.Close(); err != nil {
		logger.Println(err)
	}

	logger.Println("Server closed")
}
