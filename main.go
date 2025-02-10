package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Aloe-Corporation/logs"
	"github.com/CamilleLange/XM-company/internal/app"
	"github.com/CamilleLange/XM-company/internal/config"
)

const (
	PREFIX_ENV          = "XM_COMPANY"
	ENV_CONFIG          = PREFIX_ENV + "_CONFIG"
	DEFAULT_PATH_CONFIG = "./config/"
)

var (
	log = logs.Get()
)

func main() {
	configFilePath, present := os.LookupEnv(ENV_CONFIG)
	if !present {
		configFilePath = DEFAULT_PATH_CONFIG
	}

	config, err := config.Load(configFilePath, PREFIX_ENV)
	if err != nil {
		panic(err)
	}

	run, close, err := app.Launch(*config)
	if err != nil {
		panic(err)
	}

	go run()

	WaitSignalShutdown(close)
}

func WaitSignalShutdown(close func() error) {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	close()
}
