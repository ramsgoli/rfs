package dataserver

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	agentmanagerservice "github.com/ramsgoli/rfs/internal/master/agent_manager_service"
	chunkmanagerservice "github.com/ramsgoli/rfs/internal/master/chunk_manager_service"
)

const HAPPY_RESPONSE = 0x1
const SAD_RESPONSE = 0x2

const MAX_FILE_SIZE = 1024 // 1kb

type DataServer interface {
	Start()
	Stop()
}

type server struct {
	socket              net.Listener
	agentManagerService agentmanagerservice.AgentManagerService
	chunkManagerService chunkmanagerservice.ChunkManagerService
	wg                  *sync.WaitGroup
}

type DataServerOpts struct {
	Socket              net.Listener
	AgentManagerService agentmanagerservice.AgentManagerService
	ChunkManagerService chunkmanagerservice.ChunkManagerService
	Wg                  *sync.WaitGroup
}

func NewDataServer(opts *DataServerOpts) DataServer {
	return &server{
		socket:              opts.Socket,
		wg:                  opts.Wg,
		agentManagerService: opts.AgentManagerService,
		chunkManagerService: opts.ChunkManagerService,
	}
}

func (s *server) Start() {
	s.wg.Add(1)
	fmt.Println("Starting server...")

	for {
		conn, err := s.socket.Accept()
		if err != nil {
			log.Printf("error accepting conn: %+v", err)
			return
		}
		log.Printf("got a connection: %s\n", conn.RemoteAddr().String())
		go s.handleConn(conn)
	}
}

func (s *server) Stop() {
	s.socket.Close()
	s.wg.Done()
}

func (s *server) handleConn(conn net.Conn) {
	defer conn.Close()

	m, err := s.getWriteRequestMetadata(conn)
	if err != nil {
		log.Printf("[ERROR]: error getting write request metadata: %v", err)
		conn.Write([]byte{SAD_RESPONSE})
		return
	}

	if !s.canAcceptRequest(m) {
		log.Print("sad")
		conn.Write([]byte{SAD_RESPONSE})
	}

	if err != nil {
		log.Printf("error opening connection to agent: %v", err)
		return
	}
	s.streamClientDataToAgent(conn, m)
}

func (s *server) canAcceptRequest(m *writeRequestMetadata) bool {
	// for now, check file size
	return m.fileSize < MAX_FILE_SIZE && len(s.agentManagerService.GetAgents()) >= 1
}

func (s *server) streamClientDataToAgent(
	clientConn net.Conn,
	writeMetadata *writeRequestMetadata,
) {
	// read in chunks of 1MB
	buf := make([]byte, 1024)
	var chunk uint8 = 0
	for {
		// conn.Read() reads as many bytes as are available in the TCP buffer at the time and fit in the provided buffer.
		_, err := io.ReadFull(clientConn, buf)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				log.Printf("EOF reached.")
				break
			}
			log.Panicf("read error: %v", err)
		}
		s.chunkManagerService.HandleChunk(&chunkmanagerservice.Chunk{
			FileNameBuffer: writeMetadata.fileNameBuffer,
			ChunkIndex:     chunk,
			Data:           buf,
		})

		chunk++
	}
}

