package util

import "sync"

type Instance interface {
	OnInit()
	Run(closeSig chan bool)
	OnDestory()
}
type InstanceData struct {
	I        Instance
	closeSig chan bool
	wg       sync.WaitGroup
}

func (self *InstanceData) Init(async bool) {
	self.closeSig = make(chan bool, 1)
	self.wg.Add(1)
	self.I.OnInit()
	if async == true {
		go self.Run()
	}
	self.Run()
}
func (self *InstanceData) Run() {
	self.I.Run(self.closeSig)
	self.wg.Done()
}
func (self *InstanceData) Destory() {
	self.closeSig <- true
	self.wg.Wait()
	self.I.OnDestory()
}
