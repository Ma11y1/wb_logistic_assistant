package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"wb_logistic_assistant/internal/app"
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/logger"
)

const (
	logPath    = "./logs.log"
	configPath = "./config.json"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Logf(logger.FATAL, "Main()", "Runtime error: %v", err)
			os.Exit(1)
		}
	}()

	appConfig, err := config.NewConfigFile(configPath)
	if err != nil {
		log.Fatalf("Failed load application configuration: %v", err)
	}

	if appConfig.Debug().Path() != "" {
		file, err := os.OpenFile(appConfig.Debug().Path(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		logger.AddOutput(file)
		logger.AddOutputErr(file)
	}

	if appConfig.Debug().Enabled() {
		logger.AddOutput(os.Stdout)
		logger.AddOutputErr(os.Stdout)
	}

	application := app.NewApp(appConfig)
	if err != nil {
		logger.Logf(logger.FATAL, "Main()", "Failed to create application: %v", err)
	}

	err = application.Init()
	if err != nil {
		logger.Logf(logger.FATAL, "Main()", "Failed to initialize application: %v", err)
	}

	err = application.Start()
	if err != nil {
		logger.Logf(logger.FATAL, "Main()", "Failed to start application: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	application.Stop()
}
