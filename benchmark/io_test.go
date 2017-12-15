package benchmark

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/woobest/network/socket"

	"github.com/woobest/network"
	"github.com/woobest/util"
)

const benchmarkAddress = "127.0.0.1:7201"
const benchmarkSeconds = 10

var signal *util.SignalTester

func TestIO(t *testing.T) {

	signal = util.NewSignalTester(t)

	// 超时时间为测试时间延迟一会
	signal.SetTimeout((benchmarkSeconds + 5) * time.Second)
	server()
	go client()
	signal.WaitAndExpect("recv time out", 1)
}

func server() {
	network.MyTest()
	peer := socket.NewAcceptor().Start(benchmarkAddress)

	if peer == nil {
		return
	}
	// osSignal := make(chan os.Signal, 2)
	// signal.Notify(osSignal, os.Kill, os.Interrupt)
	// <-osSignal
}

func client() {

	peer1 := socket.NewConnector().Start(benchmarkAddress)
	if peer1 == nil {
		return
	}

	for {
		s := peer1.(*socket.Connector).GetSession()
		if s == nil {
			fmt.Print("s == nil")
			time.Sleep(time.Duration(3) * time.Second)
			continue
		}
		break
		//s.Send("hello i am client")
	}
	s := peer1.(*socket.Connector).GetSession()
	if s == nil {
		fmt.Print("s == nil")
		return
	}
	msg := network.NewMsgLen(100)
	msg.Clear()
	msg.Opcode = 11
	msg.WriteString("hello")
	s.SendPacket(msg)
	osSignal := make(chan os.Signal, 2)
	//signal.Notify(osSignal, os.Kill, os.Interrupt)
	<-osSignal
}
