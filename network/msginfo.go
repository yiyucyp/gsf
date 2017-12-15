package network

import (
	"bufio"
	"fmt"
	"io"

	"github.com/woobest/network/config"
	"github.com/woobest/util"
)

//var ByteOrder = binary.LittleEndian
var ErrorMsgInfo = fmt.Errorf("MsgInfo error")

type msgHeadStruct struct {
	Opcode   uint32
	BodySize uint16
}

const headSize = 6

type MsgInfo struct {
	msgHeadStruct
	util.Stream
}

// 利用send内存，减少申请释放
func (p *MsgInfo) BuildSendBuf(send []byte) ([]byte, error) {

	p.BodySize = uint16(len(p.Buf))
	if p.BodySize > config.MaxPacketSize || int(p.BodySize+p.HeadSize()) > cap(send) {
		return nil, ErrorMsgInfo
	}

	p.PutUint32(send, p.Opcode)
	p.PutUint16(send[4:], p.BodySize)
	if p.BodySize > 0 {
		copy(send[headSize:], p.Buf)
	}
	return send[0 : p.BodySize+headSize], nil
}

// 每包申请内存，因为可能会有把包丢到其他线程的可能
func (p *MsgInfo) BuildRecvPacket(reader *bufio.Reader) error {
	headAry := [headSize]byte{}
	_, err := io.ReadFull(reader, headAry[0:])
	if err != nil {
		return err
	}
	p.Opcode = p.Uint32(headAry[0:])
	p.BodySize = p.Uint16(headAry[4:])
	if p.BodySize > config.MaxPacketSize {
		return ErrorMsgInfo
	}
	if p.BodySize > 0 {
		p.Buf = make([]byte, p.BodySize, p.BodySize)
		_, err = io.ReadFull(reader, p.Buf)
	}
	return err
}

func (p *MsgInfo) HeadSize() uint16 {
	return headSize
}

func NewReadMsg() *MsgInfo {
	v := &MsgInfo{}
	v.ByteOrder = config.ByteOrder
	return v
}

func NewMsgLen(size int) *MsgInfo {

	v := NewReadMsg()
	v.Buf = make([]byte, size, size)
	return v
}
func MyTest() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r)
		}
	}()
	//oo.Uint16(oo.DataPtr)
}
