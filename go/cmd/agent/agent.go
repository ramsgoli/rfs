package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ramsgoli/rfs/internal/agent"
	dataserver "github.com/ramsgoli/rfs/internal/agent/data_server"
	masterservice "github.com/ramsgoli/rfs/internal/agent/master_service"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Panicf("error parsing config: %v", err)
	}

	// create data dir if doesn't exist
	if _, err := os.Stat(cfg.DataDir); os.IsNotExist(err) {
		err := os.Mkdir(cfg.DataDir, os.ModeDir)
		if err != nil {
			log.Panicf("error creating data dir: %v", err)
		}
	}

	httpClient := masterservice.NewHttpClient()
	masterService := masterservice.NewMasterService(httpClient)
	dataserver, err := dataserver.NewDataServer(&dataserver.DataServerOpts{
		Port: int64(cfg.DataPort),
		DataDir: cfg.DataDir,
	})
	if err != nil {
		log.Panicf("error starting data server: %v", err)
	}

	agent := agent.NewAgent(&agent.AgentOpts{
		MasterServerHostname: cfg.ServerHostname,
		MasterServerPort:     int64(cfg.ServerHttpPort),
		MasterService:        masterService,
		DataServer:           dataserver,
		IpAddress:            cfg.IpAddress,
	})
	err = agent.Start()
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// stop agent
	agent.Stop()
}
