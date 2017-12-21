package network

import (
	"fmt"

	"github.com/woobest/network/config"
)

const PacketHeadSize = 6

var ErrorPacketError = fmt.Errorf("Packet error")

type NetPacket struct {
	BodySize uint16
	Opcode   uint32
}

func (self *NetPacket) PrasePacket(buf []byte) error {
	self.BodySize = config.ByteOrder.Uint16(buf)
	self.Opcode = config.ByteOrder.Uint32(buf[2:])
	if self.BodySize > config.MaxPacketSize {
		return ErrorPacketError
	}
	return nil
}

func (self *NetPacket) BuildBuff(buf []byte) {
	//ary := [PacketHeadSize]byte{}
	config.ByteOrder.PutUint16(buf, self.BodySize)
	config.ByteOrder.PutUint32(buf[2:], self.Opcode)
}
