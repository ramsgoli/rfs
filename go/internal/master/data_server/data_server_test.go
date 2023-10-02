package dataserver

import (
	"net"
	"testing"
)

type fakeListener struct {
	net.Conn
}

func (f fakeListener) Accept() (net.Conn, error) {
	return f.Conn, nil
}
func (f fakeListener) Close() error {
	return nil
}
func (f fakeListener) Addr() net.Addr {
	return nil
}

func TestDataServerStart(t *testing.T) {
	client, server := net.Pipe()
	socket := fakeListener{server}

	dataServer := NewDataServer(&DataServerOpts{
		Socket: socket,
	})
}
