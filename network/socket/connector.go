package socket

import (
	"net"
	"time"

	"github.com/woobest/network"
)

// client
type Connector struct {
	stopped          bool
	address          string
	session          network.Session
	autoReconnectSec int
	network.Protocol
	network.PacketHandler
}

func (self *Connector) Start(address string) network.Peer {
	self.address = address
	go self.autoReconnectLoop()
	return self
}

func (self *Connector) Stop() {
	if self.stopped {
		return
	}
	self.stopped = true
	if self.session != nil {
		self.session.Close()
		self.session = nil
	}

}

func (self *Connector) GetSession() network.Session {
	return self.session
}

func NewConnector(protocol network.Protocol, hander network.PacketHandler) network.Peer {
	return &Connector{
		stopped:       false,
		Protocol:      protocol,
		PacketHandler: hander,
	}
}

func (self *Connector) connect() error {
	conn, err := net.Dial("tcp", self.address)
	if err != nil {
		return err
	}
	ses := newSession(conn, self)

	ses.OnClose = func(s network.Session) {
		if self.autoReconnectSec >= 0 {
			go self.autoReconnectLoop()
		}
	}
	self.session = ses
	ses.start()
	return nil
}

func (self *Connector) autoReconnectLoop() {
	for {
		err := self.connect()
		if err == nil {
			break
		}
		if self.autoReconnectSec == 0 || self.stopped {
			break
		}

		time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)
	}

}
