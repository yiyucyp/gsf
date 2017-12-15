package main

import (
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/woobest/network"
	"github.com/woobest/network/socket"
)

func main() {
	f, _ := os.Create("profile_file")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	peer := socket.NewAcceptor(network.NewProtocol(), network.NewPacketHandlerImp()).Start("127.0.0.1:10086")

	if peer == nil {
		return
	}
	osSignal := make(chan os.Signal, 2)
	signal.Notify(osSignal, os.Kill, os.Interrupt)
	<-osSignal
	peer.Stop()
}
