package network

import (
	"bufio"
	"io"
	"reflect"

	"github.com/woobest/network/config"
)

type Protocol interface {
	//完成具体的头部组织并发送，利用reflect获取msgIDgo
	SendPacket(interface{}, func([]byte) error) error
	RecvPacket([]byte, *bufio.Reader) (interface{}, error)
}

type ProtocolImp struct {
}

func (self *ProtocolImp) SendPacket(data interface{}, sendFun func([]byte) error) (err error) {
	meta, err := MessageMetaByType(reflect.TypeOf(data))
	if err != nil {
		return err
	}

	bodyBuf, err := meta.Codec.Encode(data)
	if err != nil {
		return
	}
	bodySize := len(bodyBuf)
	if bodySize > int(config.MaxPacketSize) {
		err = ErrorPacketError
		return
	}

	// send head
	head := NetPacket{}
	head.Opcode = meta.ID
	head.BodySize = uint16(bodySize)
	headArray := [PacketHeadSize]byte{}
	head.BuildBuff(headArray[0:])
	err = sendFun(headArray[0:])
	if err != nil {
		return
	}
	// send body
	err = sendFun(bodyBuf)
	return err
}

func (self *ProtocolImp) RecvPacket(buf []byte, reader *bufio.Reader) (msg interface{}, err error) {
	head := NetPacket{}
	if _, err = io.ReadFull(reader, buf[0:PacketHeadSize]); err != nil {
		return
	}
	if err = head.PrasePacket(buf); err != nil {
		return
	}

	meta, err := MessageMetaByID(head.Opcode)
	if err != nil {
		return
	}

	if _, err = io.ReadFull(reader, buf[:head.BodySize]); err != nil {
		return
	}
	msg = reflect.New(meta.Type).Interface()
	if err = meta.Codec.Decode(buf[:head.BodySize], msg); err != nil {
		return
	}
	return
}
func NewProtocol() Protocol {
	return &ProtocolImp{}
}
