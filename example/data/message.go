package data

import (
	"encoding/binary"
	"io"
)

func Read(r io.Reader) (*StringMsg,error) {
	msg := StringMsg{}
	err := binary.Read(r, binary.BigEndian, &msg.len)
	if err!=nil {
		return nil,err
	}
	msg.data = make([]byte,msg.len)
	err = binary.Read(r, binary.BigEndian, msg.data)
	if err!=nil {
		return nil,err
	}
	return &msg,nil
}

type StringMsg struct {
	len uint64
	data []byte
}

func (s *StringMsg) String() string {
	return string(s.data)
}

func New(content string) *StringMsg {
	s := StringMsg{}
	s.data = []byte(content)
	s.len = uint64(len(s.data))
	return &s
}

func (s *StringMsg) Bytes() []byte {
	buffer := make([]byte,s.len+8)
	binary.BigEndian.PutUint64(buffer[:8],s.len)
	copy(buffer[8:],s.data)
	return buffer
}
