package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var ErrorStream = fmt.Errorf("Stream out of buff")

type Stream struct {
	Buf  []byte
	Rpos int
	//bytes.Buffer
	binary.ByteOrder
}

func (p *Stream) Size() int { return len(p.Buf) }

func (p *Stream) Data() []byte { return p.Buf }

func (p *Stream) Left() int { return len(p.Buf) - p.Rpos }

func (p *Stream) Reset(buff []byte) { p.Rpos = 0; p.Buf = buff }
func (p *Stream) Clear() {
	p.Rpos = 0
	p.Buf = p.Buf[:0]
}
func (p *Stream) ReadByte() (byte, error) {
	if p.Left() < 1 {
		return 0, ErrorStream
	}
	v := p.Buf[p.Rpos]
	p.Rpos++
	return v, nil
}

func (p *Stream) ReadUint16() (uint16, error) {
	if p.Left() < 2 {
		return 0, ErrorStream
	}
	v := p.Uint16(p.Buf[p.Rpos:])
	p.Rpos += 2
	return v, nil
}

func (p *Stream) ReadUint32() (uint32, error) {
	if p.Left() < 4 {
		return 0, ErrorStream
	}
	v := p.Uint32(p.Buf[p.Rpos:])
	p.Rpos += 4
	return v, nil
}
func (p *Stream) ReadUint64() (uint64, error) {
	if p.Left() < 4 {
		return 0, ErrorStream
	}
	v := p.Uint64(p.Buf[p.Rpos:])
	p.Rpos += 8
	return v, nil
}
func (p *Stream) ReadBuff(size int) (buff []byte, err error) {
	if p.Left() < size {
		return nil, ErrorStream
	}
	buff = make([]byte, size, size)
	copy(buff, p.Buf[p.Rpos:p.Rpos+size])
	return buff, nil
}
func (p *Stream) WriteByte(v byte) {
	p.Buf = append(p.Buf, v)
}

func (p *Stream) WriteUint16(v uint16) {
	//p.PutUint16()
	t := [2]byte{}
	p.PutUint16(t[0:], v)
	p.Buf = append(p.Buf, t[0:]...)
}
func (p *Stream) WriteUint32(v uint32) {
	//p.PutUint16()
	t := [4]byte{}
	p.PutUint32(t[0:], v)
	p.Buf = append(p.Buf, t[0:]...)
}
func (p *Stream) WriteUint64(v uint64) {
	//p.PutUint16()
	t := [8]byte{}
	p.PutUint64(t[0:], v)
	p.Buf = append(p.Buf, t[0:]...)
}
func (p *Stream) WriteBuf(buf []byte) {
	p.Buf = append(p.Buf, buf...)
}

// 读结构体, slice
func (p *Stream) ReadStruct(data interface{}) error {
	dataSize := binary.Size(data)
	if dataSize == -1 {
		return ErrorStream
	}
	err := binary.Read(bytes.NewBuffer(p.Buf), p.ByteOrder, data)
	if err != nil {
		return err
	}

	p.Rpos += dataSize
	return nil
}

//  写结构体, 固定大小SLICE
func (p *Stream) WriteStruct(data interface{}) error {
	return binary.Write(bytes.NewBuffer(p.Buf), p.ByteOrder, data)
}

func (p *Stream) ReadString() (string, error) {
	v, err := p.ReadUint16()
	if err != nil || p.Left() < int(v) {
		return "", ErrorStream
	}
	bytes := p.Buf[p.Rpos : p.Rpos+int(v)]
	return string(bytes), nil
}

func (p *Stream) WriteString(str string) {
	v := []byte(str)
	p.WriteUint16(uint16(len(v)))
	p.Buf = append(p.Buf, v...)
}

func NewStream(data []byte, order binary.ByteOrder) *Stream {
	return &Stream{Buf: data, ByteOrder: order}
}
