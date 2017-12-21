package socket

import (
	"bufio"
	"io"
	"net"
	"sync/atomic"

	"github.com/woobest/network"
	"github.com/woobest/network/config"
)

const stateInit int32 = -1
const stateConnected int32 = 0
const stateCloseing int32 = 1
const stateClosed int32 = 2
const sendChanCap int32 = 32

type socketSession struct {
	OnClose     func(network.Session)
	id          int64
	p           network.Peer
	state       int32
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	writeChan   chan interface{}
	stoppedChan chan struct{}
}

func newSession(conn net.Conn, p network.Peer) *socketSession {
	return &socketSession{
		state:       -1,
		p:           p,
		conn:        conn,
		reader:      bufio.NewReader(conn),
		writer:      bufio.NewWriter(conn),
		writeChan:   make(chan interface{}, sendChanCap),
		stoppedChan: make(chan struct{}),
	}
}
func (self *socketSession) ID() int64 {
	return self.id
}
func (self *socketSession) SetID(id int64) {
	self.id = id
}
func (self *socketSession) RawConn() net.Conn { return self.conn }

func (self *socketSession) FromPeer() network.Peer { return self.p }

func (self *socketSession) Send(msg interface{}) {
	select {
	case self.writeChan <- msg:
	case <-self.stoppedChan:
		log.Errorf("Send: writeChan is closed")
	default:
		log.Errorf("Send: writeChan is full")
		self.Close()
	}
}

func (self *socketSession) start() {
	if atomic.CompareAndSwapInt32(&self.state, stateInit, stateConnected) {
		go self.sendLoop()
		go self.recvLoop()
	}
}

func (self *socketSession) Close() {
	if atomic.CompareAndSwapInt32(&self.state, stateConnected, stateCloseing) {
		//self.conn.(*net.TCPConn).SetLinger(0)
		self.conn.Close()       // 退出recvLoop
		close(self.stoppedChan) // 退出sendLoop
		if self.OnClose != nil {
			self.OnClose(self) // 要换考虑多线程问题
		}
		log.Infof("Close")
	} else {
		log.Errorf("Close state=%d", self.state)
	}
}
func (self *socketSession) sendLoop() {
	defer self.Close()
	peer := self.FromPeer()
	for {
		select {
		case msg, ok := <-self.writeChan:
			{
				if msg == nil || ok == false { // nil主动关闭,或者chan关掉了退出并回收
					return
				}
				err := peer.SendPacket(msg, func(sendbuf []byte) error {
					return self.writeBuf(sendbuf)
				})
				if err != nil {
					return
				}
			}
		case <-self.stoppedChan:
			return
		default:
			{
				var err error
				for i := 0; i < 100; i++ {
					if err = self.writer.Flush(); err != io.ErrShortWrite {
						//log.Infof("send default")
						break
					}
					if err != nil {
						return
					}
				}
				msg, ok := <-self.writeChan
				if msg == nil || ok == false { // nil主动关闭,或者chan关掉了退出并回收
					return
				}
				err = peer.SendPacket(msg, func(sendbuf []byte) error {
					return self.writeBuf(sendbuf)
				})
				if err != nil {
					return
				}

			}
		}
	}
}
func (self *socketSession) writeBuf(buff []byte) (err error) {
	length := len(buff)
	var n, nn int
	for n < length && err == nil {
		nn, err = self.writer.Write(buff[n:])
		n += nn
	}
	if err != nil {
		log.Errorf("writeBuf write error:%s", err.Error())
	}
	return
}

func (self *socketSession) recvLoop() {
	defer self.Close()
	peer := self.FromPeer()
	// 减少slice的申请释放
	buf := make([]byte, config.MaxPacketSize, config.MaxPacketSize)
	for {
		_, err := peer.RecvPacket(buf, func(buf []byte) error {
			_, err := io.ReadFull(self.reader, buf)
			return err
		}, self)
		if err != nil {
			log.Errorf("recvLoop error:%s", err.Error())
			break
		}
	}
}
