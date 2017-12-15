package config

import "encoding/binary"

var (
	ByteOrder     binary.ByteOrder = binary.LittleEndian
	MaxPacketSize uint16           = 10240
	LenStackBuf                    = 4096
)

func Init() {

}
