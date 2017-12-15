package json

import (
	"encoding/json"

	"github.com/woobest/network"
)

type jsonCodec struct {
}

func (self *jsonCodec) Name() string {
	return "json"
}

func (self *jsonCodec) Encode(msgObj interface{}) ([]byte, error) {

	return json.Marshal(msgObj)

}

func (self *jsonCodec) Decode(data []byte, msgObj interface{}) error {

	return json.Unmarshal(data, msgObj)
}

func init() {

	network.RegisterCodec("json", new(jsonCodec))
}
