package agentserver

import (
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	agentmanagerservice "github.com/ramsgoli/rfs/internal/master/agent_manager_service"
)

type agentServer struct {
	port uint
	g    *gin.Engine
	wg   *sync.WaitGroup

	agentManagerService agentmanagerservice.AgentManagerService
}

type AgentServerOpts struct {
	Port                uint
	Wg                  *sync.WaitGroup
	AgentManagerService agentmanagerservice.AgentManagerService
}

func NewAgentServer(opts *AgentServerOpts) AgentServer {
	r := gin.Default()
	return &agentServer{
		port:                opts.Port,
		wg:                  opts.Wg,
		g:                   r,
		agentManagerService: opts.AgentManagerService,
	}
}

type AgentServer interface {
	Start()
	Stop()
}

type AgentRequest struct {
	ClientIP string `json:"ip_address"`
	Port     int64  `json:"port"`
}

func (s *agentServer) Start() {
	s.wg.Add(1)

	s.g.POST("/node/register", s.nodeRegisterRoute())
	s.g.Run(fmt.Sprintf(":%d", s.port))
}

func (s *agentServer) Stop() {
	s.wg.Done()
}

func (s *agentServer) nodeRegisterRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req AgentRequest
		if err := c.BindJSON(&req); err != nil {
			c.Writer.Write([]byte(`{"success": false}`))
			return
		}

		if err := s.agentManagerService.RegisterAgent(&agentmanagerservice.AgentDetails{
			ClientIP: req.ClientIP,
			Port:     req.Port,
		}); err != nil {
			log.Printf("error registering agent: %v\n", err)
			c.Writer.Write([]byte(`{"success": false}`))
			return
		}

		c.Writer.Write([]byte(`{"success": true}`))
	}
}
