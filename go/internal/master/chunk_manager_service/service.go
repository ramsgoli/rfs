package chunkmanagerservice

import (
	"fmt"
	"math"

	agentmanagerservice "github.com/ramsgoli/rfs/internal/master/agent_manager_service"
)

type ChunkManagerService interface {
	HandleChunk(*Chunk) error
}

type Chunk struct {
	FileNameBuffer [16]byte
	ChunkIndex uint8
	Data []byte
}

type chunkManagerService struct {
	agentManagerService agentmanagerservice.AgentManagerService
	chunkMetadata map[*Chunk]*agentmanagerservice.AgentDetails
}

type ChunkManagerServiceOpts struct {
	AgentManagerService agentmanagerservice.AgentManagerService
}

func NewChunkManagerService(o *ChunkManagerServiceOpts) *chunkManagerService {
	return &chunkManagerService{
		agentManagerService: o.AgentManagerService,
		chunkMetadata: make(map[*Chunk]*agentmanagerservice.AgentDetails),
	}
}

func (c *chunkManagerService) HandleChunk(chunk *Chunk) error {
	agent := c.getAgentToWriteNextChunk()
	agentConn, err := c.agentManagerService.GetAgentConnection(agent)
	if err != nil {
		return fmt.Errorf("could not write chunk: %v", err)
	}

	defer agentConn.Close()

	dataSize := 16 + 1 +len(chunk.Data)
	buf := make([]byte, dataSize)

	// prepare the buffer to send to agent
	copy(buf[:], chunk.FileNameBuffer[:])
	buf[16] = chunk.ChunkIndex
	copy(buf[17:], chunk.Data)

	agentConn.WriteData(buf)
	c.recordChunkWrittenToAgent(chunk, agent)

	return nil
}

func (c *chunkManagerService) getAgentToWriteNextChunk() (*agentmanagerservice.AgentDetails) {
	allAgents := c.agentManagerService.GetAgents()
	agentToChunkCount := make(map[*agentmanagerservice.AgentDetails]int)
	for _, a := range allAgents {
		agentToChunkCount[a] = 0
	}

	for _, agent := range c.chunkMetadata {
		agentToChunkCount[agent] += 1
	}

	// find the agent with the least number of chunks
	var agent *agentmanagerservice.AgentDetails
	minChunkCount := math.Inf(1)
	for a, chunkCount := range agentToChunkCount {
		if float64(chunkCount) < minChunkCount {
			agent = a
			minChunkCount = float64(chunkCount)
		}
	}
	return agent
}

func (c *chunkManagerService) recordChunkWrittenToAgent(chunk *Chunk, agent *agentmanagerservice.AgentDetails) {
	c.chunkMetadata[chunk] = agent
}
