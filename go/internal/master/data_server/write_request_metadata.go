package dataserver

import (
	"encoding/binary"
	"fmt"
	"net"
)

type writeRequestMetadata struct {
	fileNameBuffer [16]byte
	fileName       string
	fileSizeBuffer [4]byte
	fileSize       uint32
}

func (s *server) getWriteRequestMetadata(conn net.Conn) (*writeRequestMetadata, error) {
	magicByteBuffer := make([]byte, 1)
	_, err := conn.Read(magicByteBuffer)
	if err != nil {
		return nil, fmt.Errorf("error reading magic byte: %v", err)
	}
	if string(magicByteBuffer[0]) != "x" {
		return nil, fmt.Errorf("did not find magic byte")
	}

	// read fileName length
	var fileNameBuffer [16]byte
	_, err = conn.Read(fileNameBuffer[:])
	if err != nil {
		return nil, fmt.Errorf("error reading file name: %v", err)
	}
	fileName := string(fileNameBuffer[:])

	// read file data length
	var fileSizeBuffer [4]byte
	_, err = conn.Read(fileSizeBuffer[:])
	if err != nil {
		return nil, fmt.Errorf("error reading file length: %v", err)
	}
	fileSize := binary.BigEndian.Uint32(fileSizeBuffer[:])

	return &writeRequestMetadata{
		fileNameBuffer: fileNameBuffer,
		fileName:       fileName,
		fileSizeBuffer: fileSizeBuffer,
		fileSize:       fileSize,
	}, nil
}
