package main

import (
	"os"
	"os/signal"
	"reflect"
	"runtime/pprof"

	"github.com/woobest/network"
	"github.com/woobest/network/socket"
	"github.com/woobest/protocol/pb/msgdef"
)

func main() {
	f, _ := os.Create("profile_file")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	network.RegisterMessageMeta("pb", 1, reflect.TypeOf((*msgdef.TestEchoACK)(nil)).Elem(), func(s network.Session, msg interface{}, meta *network.MessageMeta) {
		pack := msg.(*msgdef.TestEchoACK)
		pack.Content = "hi"
		//fmt.Println(pack.Content)
	})

	peer := socket.NewAcceptor(network.NewProtocol()).Start("127.0.0.1:10086")

	if peer == nil {
		return
	}
	osSignal := make(chan os.Signal, 2)
	signal.Notify(osSignal, os.Kill, os.Interrupt)
	<-osSignal
	peer.Stop()
}
