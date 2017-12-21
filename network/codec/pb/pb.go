package pb

import (
	"github.com/golang/protobuf/proto"
	"github.com/woobest/network/codec"
)

type pbCodec struct {
}

func (self *pbCodec) Name() string {
	return "pb"
}

func (self *pbCodec) Encode(msgObj interface{}) ([]byte, error) {

	msg := msgObj.(proto.Message)

	data, err := proto.Marshal(msg)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (self *pbCodec) Decode(data []byte, msgObj interface{}) error {

	err := proto.Unmarshal(data, msgObj.(proto.Message))

	if err != nil {
		return err
	}

	return nil
}

func init() {

	codec.RegisterCodec("pb", new(pbCodec))
}
