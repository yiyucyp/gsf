package benchmark

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/woobest/network"
	_ "github.com/woobest/network/codec/pb"
	"github.com/woobest/network/socket"
	"github.com/woobest/protocol/pb/msgdef"
	"github.com/woobest/util"
)

const benchmarkAddress = "127.0.0.1:7201"
const benchmarkSeconds = 10
const benchmarkClients = 100

var signal *util.SignalTester
var peerServer network.Peer

func TestIO(t *testing.T) {

	signal = util.NewSignalTester(t)

	// 超时时间为测试时间延迟一会
	signal.SetTimeout((benchmarkSeconds + 5) * time.Second)
	fmt.Print("server start")
	server()
	fmt.Print("server end")
	for i := 0; i < benchmarkClients; i++ {
		go client()
	}
	fmt.Print("client end")

	signal.WaitAndExpect("recv time out", 1)
}

func server() {
	qpsm := NewQPSMeter(func(qps int) {

		fmt.Print("QPS: %d", qps)

	})

	network.RegisterMessageMeta("pb", 1, reflect.TypeOf((*msgdef.TestEchoACK)(nil)).Elem(), func(s network.Session, msg interface{}, meta *network.MessageMeta) {
		pack := msg.(*msgdef.TestEchoACK)
		pack.Content = "hi"
		//fmt.Println(pack.Content)
		if qpsm.Acc() > benchmarkSeconds {
			fmt.Print("Average QPS: %d", qpsm.Average())
			signal.Done(1)
			//fmt.Print("Average QPS: %d", qpsm.Average())
		}
	})
	peerServer := socket.NewAcceptor(network.NewProtocol())

	if peerServer == nil {
		return
	}
	peerServer.Start("127.0.0.1:10086")
	//signal.WaitAndExpect("recv time out", 2)
}

type Client struct {
	client network.Peer
}

func (self *Client) OnInit() {

	self.client = socket.NewConnector(network.NewProtocol())
	network.RegisterMessageMeta("pb", 1, reflect.TypeOf((*msgdef.TestEchoACK)(nil)).Elem(), func(s network.Session, msg interface{}, meta *network.MessageMeta) {
		//pack := msg.(*msgdef.TestEchoACK)
		//fmt.Println(pack.Content)
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
					//fmt.Print("send")
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
func client() {

	var ins = util.InstanceData{I: &Client{}}
	ins.Init(false)
	ins.Destory()
}
