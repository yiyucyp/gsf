package benchmark

import (
	"sync"
	"time"
)

type QPSMeter struct {
	qpsGuard sync.Mutex
	qps      int
	total    int

	count int
}

func (self *QPSMeter) Acc() int {
	self.qpsGuard.Lock()

	defer self.qpsGuard.Unlock()

	self.qps++

	return self.count
}

// 一轮计算
func (self *QPSMeter) Turn() (ret int) {
	self.qpsGuard.Lock()

	if self.qps > 0 {
		ret = self.qps
	}

	self.total += self.qps

	self.qps = 0
	self.count++

	self.qpsGuard.Unlock()

	return
}

// 均值
func (self *QPSMeter) Average() int {

	self.qpsGuard.Lock()

	defer self.qpsGuard.Unlock()

	if self.count == 0 {
		return 0
	}

	return self.total / self.count
}

func NewQPSMeter(callback func(int)) *QPSMeter {

	self := &QPSMeter{}

	time.AfterFunc(time.Second*time.Duration(10), func() {
		self.count = 11
		callback(self.qps)
	})
	// timer.NewLoop(time.Second, func(ctx *timer.Loop) {

	// 	qps := self.Turn()

	// 	callback(qps)

	// }, nil).Start()

	return self
}
