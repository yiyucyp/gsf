package network

import (
	"fmt"
	"reflect"
)

type MessageMeta struct {
	Type  reflect.Type
	ID    uint32
	Codec Codec
}

var (
	metaByID          = map[uint32]*MessageMeta{}
	metaByType        = map[reflect.Type]*MessageMeta{}
	ErrorsMessageMeta = fmt.Errorf("MessageMeta error")
)

func RegisterMessageMeta(msgID uint32, msgType reflect.Type) {
	meta := &MessageMeta{ID: msgID, Type: msgType}
	if msgType.Kind() == reflect.Ptr {
		msgType = msgType.Elem()
	}
	if _, ok := metaByID[msgID]; ok {
		panic(fmt.Sprintf("dumplicate message meta id:%d", msgID))
	}
	if _, ok := metaByType[msgType]; ok {
		panic("dumplicate message meta type:" + msgType.String())
	}
	metaByID[msgID] = meta
	metaByType[msgType] = meta
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
