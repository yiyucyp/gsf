package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/woobest/network"
	"github.com/woobest/network/socket"
	"github.com/woobest/util"
)

type Client struct {
	client network.Peer
}

func (self *Client) OnInit() {
	self.client = socket.NewConnector(network.NewProtocol(), network.NewPacketHandlerImp())
	self.client.Start("127.0.0.1:10086")
}

func (self *Client) Run(closeSing chan bool) {
	for {
		select {
		case <-closeSing:
			return
		}

		time.Sleep(time.Duration(10) * time.Millisecond)
		s := self.client.(*socket.Connector).GetSession()
		if s != nil {
			//send
		}
	}
}
func (self *Client) OnDestory() {
	self.client.Stop()
}
func main() {
	var ins = util.InstanceData{I: &Client{}}
	ins.Init()
	// msg := network.NewMsgLen(100)
	// msg.Clear()
	// msg.Opcode = 11
	// msg.WriteString("hello")
	// s.SendPacket(msg)
	osSignal := make(chan os.Signal, 2)
	signal.Notify(osSignal, os.Kill, os.Interrupt)
	<-osSignal
	ins.Destory()
}
