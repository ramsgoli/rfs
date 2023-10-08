package agent

import (
	"context"
	"log"

	dataserver "github.com/ramsgoli/rfs/internal/agent/data_server"
	masterservice "github.com/ramsgoli/rfs/internal/agent/master_service"
)

type agent struct {
	id                    int64
	masterServiceHostname string
	masterServicePort     int64
	masterService         masterservice.MasterService
	dataServer            dataserver.DataServer
	ipAddress             string
	cancelCtx             context.Context
}

type AgentOpts struct {
	Id                   int64
	MasterServerHostname string
	MasterServerPort     int64
	MasterService        masterservice.MasterService
	DataServer           dataserver.DataServer
	IpAddress            string
}

type Agent interface {
	Start() error
	Stop()
}

func NewAgent(opts *AgentOpts) Agent {
	return &agent{
		id:                    opts.Id,
		masterServiceHostname: opts.MasterServerHostname,
		masterServicePort:     opts.MasterServerPort,
		masterService:         opts.MasterService,
		ipAddress:             opts.IpAddress,
		dataServer:            opts.DataServer,
	}
}

func (a *agent) Start() error {
	go func() {
		err := a.dataServer.OpenSocket()
		if err != nil {
			log.Printf("error opening socket: %v", err)
		}
	}()

	err := a.masterService.RegisterWithMaster(
		a.masterServiceHostname,
		a.masterServicePort,
		a.ipAddress,
		a.dataServer.GetPort(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *agent) Stop() {
	a.dataServer.CloseSocket()
}
