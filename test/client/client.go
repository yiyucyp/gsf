package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/woobest/network"
	"github.com/woobest/network/socket"
	"github.com/woobest/protocol/pb/msgdef"
	"github.com/woobest/util"
)

type Client struct {
	client network.Peer
}

func (self *Client) OnInit() {
	self.client = socket.NewConnector(network.NewProtocol())
	network.RegisterMessageMeta("pb", 1, reflect.TypeOf(msgdef.TestEchoACK{}), func(s network.Session, msg interface{}, meta *network.MessageMeta) {
		pack := msg.(*msgdef.TestEchoACK)
		fmt.Println(pack.Content)
	})
	self.client.Start("127.0.0.1:10086")
}

func (self *Client) Run(closeSing chan bool) {
	for {
		select {
		case <-closeSing:
		default:
			{
				time.Sleep(time.Duration(10) * time.Millisecond)
				s := self.client.(*socket.Connector).GetSession()
				if s != nil {
					//send
					//msg := msgdef.TestEchoACK{Content: "hello", Bytes: make([]byte, 10, 10)}
					fmt.Print("send")
					s.Send(&msgdef.TestEchoACK{Content: "hello", Bytes: make([]byte, 10, 10)})
				}
			}
			//return
		}
	}
}
func (self *Client) OnDestory() {
	self.client.Stop()
}
func main() {
	var ins = util.InstanceData{I: &Client{}}
	ins.Init(true)
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
