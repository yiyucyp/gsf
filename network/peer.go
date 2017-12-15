package network

import (
	"net"
)

type Session interface {
	Send(msg interface{})
	Close()
	ID() int64
	RawConn() net.Conn
	FromPeer() Peer
}

type Peer interface {
	Start(address string) Peer
	Stop()
	Protocol
	PacketHandler
}
