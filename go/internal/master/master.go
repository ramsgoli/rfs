package master

import (
	"log"

	agentmanagerservice "github.com/ramsgoli/rfs/internal/master/agent_manager_service"
	agentserver "github.com/ramsgoli/rfs/internal/master/agent_server"
	dataserver "github.com/ramsgoli/rfs/internal/master/data_server"
)

type master struct {
	agentServer         agentserver.AgentServer
	dataServer          dataserver.DataServer
	agentManagerService agentmanagerservice.AgentManagerService
}

type MasterService interface {
	Start()
	Stop()
}

func NewMasterService(
	agentServer agentserver.AgentServer,
	dataServer dataserver.DataServer,
	agentManagerService agentmanagerservice.AgentManagerService,
) MasterService {
	return &master{
		agentServer:         agentServer,
		dataServer:          dataServer,
		agentManagerService: agentManagerService,
	}
}

func (m *master) Start() {
	go m.agentServer.Start()
	go m.dataServer.Start()
}
func (m *master) Stop() {
	m.agentServer.Stop()
	m.dataServer.Stop()

	log.Printf("I ended with the following registered agents: %+v\n", m.agentManagerService.GetAgents())
}
