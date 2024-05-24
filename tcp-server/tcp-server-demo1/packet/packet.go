package packet

import (
	"bytes"
	"fmt"
)

// Packet协议定义

/*
### packet header
1 byte: commandID

### submit packet

8字节 ID 字符串
任意字节 payload

### submit ack packet

8字节 ID 字符串
1字节 result
*/

const (
	CommandConn   = iota + 0x01 // 0x01
	CommandSubmit               // 0x02
)

const (
	CommandConnAck   = iota + 0x80 // 0x81
	CommandSubmitAck               //0x82
)

type Packet interface {
	Decode([]byte) error     // []byte -> struct
	Encode() ([]byte, error) //  struct -> []byte
}

type Submit struct {
	ID      string
	Payload []byte
}

func (s *Submit) Decode(pktBody []byte) error {
	s.ID = string(pktBody[:8]) // ID八字节
	s.Payload = pktBody[8:]    // 后面全是payload
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), s.Payload}, nil), nil
}

type SubmitAck struct {
	ID     string
	Result uint8
}

func (s *SubmitAck) Decode(pktBody []byte) error {
	s.ID = string(pktBody[0:8])  // 0-7赋值给S.ID
	s.Result = uint8(pktBody[8]) // uint8一个字节大小，所以第八个字节赋值给Result
	return nil
}

func (s *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), []byte{s.Result}}, nil), nil
}

// TODO:自己code conn的Decode&Encode
/*type Connect struct {
	ID      string
	Payload []byte
}

func (c *Connect) Decode(pktBody []byte) error {
	c.ID = string(pktBody[:8]) //commandID八字节
	c.Payload = pktBody[8:]    // 后面全是payload
	return nil
}

func (c *Connect) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(c.ID[:8]), c.Payload}, nil), nil
}

type ConnectAck struct {
	ID     string
	Result uint8
}

func (c *ConnectAck) Decode(pktBody []byte) error {
	c.ID = string(pktBody[0:8]) // 0-7赋值给c.ID
	c.Result = pktBody[8]       //  uint8一个字节大小，所以第八个字节赋值给Result
	return nil
}

func (c *ConnectAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(c.ID[:8]), []byte{c.Result}}, nil), nil
}*/

// Decode packet包级的Decode，将字节转为packet，两个包级的Decode&Encode提供给Frame包使用
func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]
	pktBody := packet[1:]

	switch commandID {
	// TODO: 需要自己code conn
	case CommandConn:
		return nil, nil
	case CommandConnAck:
		return nil, nil
	case CommandSubmit:
		s := Submit{}
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	case CommandSubmitAck:
		s := SubmitAck{}
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, fmt.Errorf("unknown commandID [%d]", commandID)
	}
}

// Encode 解码将packet转为byte提供给frame调用
func Encode(p Packet) ([]byte, error) {
	var commandID uint8
	var pktBody []byte
	var err error

	switch t := p.(type) {
	case *Submit:
		commandID = CommandSubmit
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{[]byte{commandID}, pktBody}, nil), nil
}
