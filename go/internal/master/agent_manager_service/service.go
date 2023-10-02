package agentmanagerservice

import (
	"fmt"
	"net"
)

type AgentManagerService interface {
	RegisterAgent(opts *AgentDetails) error
	GetAgents() []*AgentDetails
	GetAgentConnection(*AgentDetails) (*AgentConnection, error)
}

type AgentConnection struct {
	conn net.Conn
}

func (c *AgentConnection) WriteData(data []byte) {
	// write magic byte
	c.conn.Write([]byte{byte('x')})

	c.conn.Write(data)
}

func (c *AgentConnection) Close() {
	c.conn.Close()
}

func NewAgentManagerService() AgentManagerService {
	return &agentManagerService{
		registeredAgents: make([]*AgentDetails, 0),
	}
}

type AgentDetails struct {
	ClientIP string
	Port     int64
}

type agentManagerService struct {
	registeredAgents []*AgentDetails
}

func (s *agentManagerService) RegisterAgent(opts *AgentDetails) error {
	// attempt to open connection to client
	// conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", opts.ClientIP, opts.Port))
	// if err != nil {
	// 	return fmt.Errorf("error connecting to agent: %v", err)
	// }
	// conn.Close()

	s.registeredAgents = append(s.registeredAgents, opts)
	return nil
}

func (s *agentManagerService) GetAgents() []*AgentDetails {
	return s.registeredAgents
}

// simple: get a connection to the agent by IP
func (s *agentManagerService) GetAgentConnection(a *AgentDetails) (*AgentConnection, error) {
	if a == nil {
		return nil, fmt.Errorf("could not find agent")
	}

	socket, err := net.Dial("tcp", fmt.Sprintf("%s:%d", a.ClientIP, a.Port))
	if err != nil {
		return nil, fmt.Errorf("error opening connection to agent: %v", err)
	}

	connection := &AgentConnection{
		conn: socket,
	}
	return connection, nil
}
