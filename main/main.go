package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/anthill-com/ImageProcessorService/main/handler"
	"github.com/anthill-com/ImageProcessorService/main/handler/utils"
)

func main() {
	config := utils.LoadConfiguration("config.toml")

	logger, logFile := utils.CreateLog(config.LogFilePath)
	defer logFile.Close()

	dataBase := utils.CreateDB(logger, config)
	dataBase.CreateTable()

	selector := handler.CreateSelector(logger, dataBase, config, utils.CreateValidator(logger, config))

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
