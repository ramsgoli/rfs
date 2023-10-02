package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ramsgoli/go_sockets/internal/master"
	agentmanagerservice "github.com/ramsgoli/go_sockets/internal/master/agent_manager_service"
	agentserver "github.com/ramsgoli/go_sockets/internal/master/agent_server"
	chunkmanagerservice "github.com/ramsgoli/go_sockets/internal/master/chunk_manager_service"
	dataserver "github.com/ramsgoli/go_sockets/internal/master/data_server"
)

const DataServerPort = 8000
const AgentServerPort = 8080

func main() {
	var wg sync.WaitGroup

	agentManagerService := agentmanagerservice.NewAgentManagerService()
	chunkManagerService := chunkmanagerservice.NewChunkManagerService(&chunkmanagerservice.ChunkManagerServiceOpts{
		AgentManagerService: agentManagerService,
	})

	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", DataServerPort))
	if err != nil {
		log.Panicf("error opening data socket: %v", err)
	}
	dataServer := dataserver.NewDataServer(&dataserver.DataServerOpts{
		AgentManagerService: agentManagerService,
		Wg:                  &wg,
		Socket:              socket,
		ChunkManagerService: chunkManagerService,
	})
	agentServer := agentserver.NewAgentServer(&agentserver.AgentServerOpts{
		Port:                AgentServerPort,
		Wg:                  &wg,
		AgentManagerService: agentManagerService,
	})

	master := master.NewMasterService(agentServer, dataServer, agentManagerService)
	go master.Start()

	// wait for shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("received shutdown signal")

	// shut down services
	master.Stop()

	wg.Wait()
	log.Println("shutting down")
}
