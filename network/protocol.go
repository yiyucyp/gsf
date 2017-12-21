package network

import (
	"reflect"

	"github.com/woobest/network/config"
)

type ProtocolDefault struct {
}

func (self *ProtocolDefault) SendPacket(data interface{}, sendFun func([]byte) error) (err error) {
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

func (self *ProtocolDefault) RecvPacket(buf []byte, readFun func([]byte) error, s Session) (msg interface{}, err error) {
	head := NetPacket{}
	if err = readFun(buf[0:PacketHeadSize]); err != nil {
		return
	}
	if err = head.PrasePacket(buf); err != nil {
		return
	}

	meta, err := MessageMetaByID(head.Opcode)
	if err != nil {
		return
	}

	if err = readFun(buf[:head.BodySize]); err != nil {
		return
	}
	pack := reflect.New(meta.Type).Interface()
	if err = meta.Codec.Decode(buf[:head.BodySize], pack); err != nil {
		return
	}
	meta.Hander(s, pack, meta)
	return
}
func NewProtocol() Protocol {
	return &ProtocolDefault{}
}
