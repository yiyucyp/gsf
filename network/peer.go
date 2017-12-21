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
	//PacketHandler
}

type Protocol interface {
	//完成具体的头部组织并发送，利用reflect获取msgIDgo
	SendPacket(interface{}, func([]byte) error) error
	RecvPacket([]byte, func([]byte) error, Session) (interface{}, error)
}

type PacketHandler interface {
	HandlePacket(packet interface{}, s Session)
}
