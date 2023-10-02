package dataserver

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type dataServer struct {
	port   int64
	socket net.Listener
	dataDir string
}

type DataServerOpts struct {
	Port int64
	DataDir string
}

type DataServer interface {
	OpenSocket() error
	CloseSocket()
	GetPort() int64
}

func NewDataServer(opts *DataServerOpts) (DataServer, error) {
	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.Port))
	if err != nil {
		return nil, fmt.Errorf("error opening tcp socket: %v", err)
	}

	return &dataServer{
		port:   opts.Port,
		socket: socket,
		dataDir: opts.DataDir,
	}, nil
}

func (d *dataServer) GetPort() int64 {
	return d.port
}

func (d *dataServer) OpenSocket() error {
	for {
		conn, err := d.socket.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %v", err)
		}

		go d.handleConn(conn)
	}
}

func (d *dataServer) CloseSocket() {
	d.socket.Close()
}

func (d *dataServer) handleConn(conn net.Conn) {
	// magic byte + filename buffer + chunk idx + data
	bufferSize := 1 + 16 + 1 + 1024
	complete := false
	for !complete {
		buf := make([]byte, bufferSize)
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			if err != io.ErrUnexpectedEOF {
				log.Printf("error reading data: %v", err)
				break
			}
			log.Printf("EOF reached")
			complete = true
		}
		fileName := string(bytes.Trim(buf[1:17], "\x00"))
		fileChunkIndex := int(buf[17])

		f, err := os.Create(fmt.Sprintf("%s/%s_%d", d.dataDir, fileName, fileChunkIndex))
		if err != nil {
			log.Printf("error writing file: %v\n", err)
			break
		}
		f.Write(bytes.Trim(buf[18:], "\x00"))
	}
}

func printBytes(bytes []byte) {
	for _, b := range bytes {
		fmt.Printf("%c", b)
	}
}
