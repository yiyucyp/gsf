package binary

import (
	"bytes"
	"encoding/binary"

	"github.com/woobest/network"
)

type binaryCodec struct {
	binary.ByteOrder
}

func (self *binaryCodec) Name() string {
	return "binary"
}

func (self *binaryCodec) Encode(msgObj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, self.ByteOrder, msgObj)
	return buf.Bytes(), err
	//return goobjfmt.BinaryWrite(msgObj)

}

func (self *binaryCodec) Decode(data []byte, msgObj interface{}) error {

	err := binary.Write(bytes.NewBuffer(data), self.ByteOrder, msgObj)
	return err
}

func init() {

	network.RegisterCodec("binary", new(binaryCodec))
}

func NewBinaryCode(order binary.ByteOrder) network.Codec {
	return &binaryCodec{ByteOrder: order}
}
