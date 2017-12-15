package socket

import "github.com/woobest/network"
import "net"

type socketAcceptor struct {
	network.SessionManager
	stopping  chan bool
	listerner net.Listener
	address   string
	network.Protocol
	network.PacketHandler
}

func (self *socketAcceptor) waitStopFinished() {
	if self.stopping != nil {
		<-self.stopping
		self.stopping = nil
	}
}
func (self *socketAcceptor) isStopping() bool {
	return self.stopping != nil
}

func (self *socketAcceptor) startStopping() {
	self.stopping = make(chan bool)
}
func (self *socketAcceptor) endStopping() {
	select {
	case self.stopping <- true:
	default:
		self.stopping = nil
	}
}

func (self *socketAcceptor) Start(address string) network.Peer {
	self.waitStopFinished()
	if self.isStopping() {
		return self
	}
	self.address = address
	ln, err := net.Listen("tcp", address)

	self.listerner = ln
	if err != nil {

	}
	go self.accept()
	return self
}

func (self *socketAcceptor) Stop() {
	if self.isStopping() {
		return
	}
	self.startStopping()
	self.listerner.Close()
	self.CloseAllSession()
	self.waitStopFinished()
}

func (self *socketAcceptor) accept() {
	for {

		if self.isStopping() {
			break
		}
		conn, err := self.listerner.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				continue
			}
		}
		go self.onAccepted(conn)
	}
	self.endStopping()
}

func (self *socketAcceptor) onAccepted(conn net.Conn) {
	session := newSession(conn, self)
	session.start()
}

func NewAcceptor(procotol network.Protocol, hander network.PacketHandler) network.Peer {
	peer := &socketAcceptor{
		SessionManager: network.NewSessionManager(),
		Protocol:       procotol,
		PacketHandler:  hander,
	}

	return peer
}
