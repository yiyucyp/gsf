package network

import (
	"fmt"
	"reflect"

	"github.com/woobest/network/codec"
	_ "github.com/woobest/network/codec/pb" //增加序列化协议
)

type MessageHandler func(session Session, packet interface{}, meta *MessageMeta)
type MessageMeta struct {
	Type   reflect.Type
	ID     uint32
	Codec  codec.Codec
	Hander MessageHandler
}

var (
	metaByID          = map[uint32]*MessageMeta{}
	metaByType        = map[reflect.Type]*MessageMeta{}
	ErrorsMessageMeta = fmt.Errorf("MessageMeta error")
)

// 统一一个地方调用？通过工具来实现？具体的响应函数可以后续再来定义
func RegisterMessageMeta(code string, msgID uint32, msgType reflect.Type, hander MessageHandler) {

	meta := &MessageMeta{ID: msgID, Type: msgType, Hander: hander}
	if msgType.Kind() == reflect.Ptr {
		panic("dumplicate message meta is ptr:")
		msgType = msgType.Elem()
	}
	if _, ok := metaByID[msgID]; ok {
		//panic(fmt.Sprintf("dumplicate message meta id:%d", msgID))
	}
	if _, ok := metaByType[msgType]; ok {
		//panic("dumplicate message meta type:" + msgType.String())
	}
	coder := codec.FetchCodec(code)
	if coder == nil {
		panic("dumplicate message meta codec:" + code)
	}
	meta.Codec = coder
	metaByID[msgID] = meta
	metaByType[msgType] = meta

}

// 方便在各自的逻辑代码文件中直接init()里面指定响应函数
func RegisterHandler(msgType reflect.Type, handler MessageHandler) {
	if msgType.Kind() == reflect.Ptr {
		panic("dumplicate message meta is ptr:")
	}
	if _, ok := metaByType[msgType]; ok {
		metaByType[msgType].Hander = handler
	}
}
func MessageMetaByID(id uint32) (*MessageMeta, error) {
	if v, ok := metaByID[id]; ok && v.Codec != nil {
		return v, nil
	}
	return nil, ErrorsMessageMeta
}

func MessageMetaByType(t reflect.Type) (*MessageMeta, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v, ok := metaByType[t]; ok && v.Codec != nil {
		return v, nil
	}
	return nil, ErrorsMessageMeta
}
