package packet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*主要test submit， submitAck， 包级别的decode和encode*/

func TestSubmit_Decode(t *testing.T) {
	s := &Submit{}
	var b []byte
	var ID = "12345678"
	var Payload = "hello"
	b = append(b, []byte(ID)...)
	b = append(b, []byte(Payload)...)
	err := s.Decode(b)

	result := &Submit{
		ID:      ID,
		Payload: []byte(Payload),
	}
	assert.Equal(t, result, s, "The two words should be the same.")

	if err != nil {
		t.Error(err)
	}
}

func TestSubmit_Encode(t *testing.T) {
	var ID = "12345678"
	var Payload = "hello"

	s := &Submit{}
	s.ID = ID
	s.Payload = []byte(Payload)

	encode, err := s.Encode()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(encode)
}

func TestSubmitAck_Decode(t *testing.T) {
	s := &SubmitAck{}
	var b []byte
	var ID = "12345678"
	var Result = uint8(0)

	b = append(b, []byte(ID)...)
	b = append(b, Result)
	err := s.Decode(b)
	if err != nil {
		t.Error(err)
	}
}

func TestSubmitAck_Encode(t *testing.T) {
	var ID = "12345678"
	var Result = uint8(0)

	s := &SubmitAck{
		ID:     ID,
		Result: Result,
	}
	bytes, err := s.Encode()
	fmt.Println(bytes)
	if err != nil {
		t.Error(err)
	}
}

func TestDecodeSubmit(t *testing.T) {
	var b []byte
	var ID = "12345678"
	var Payload = "hello"
	b = append(b, CommandSubmit)
	b = append(b, []byte(ID)...)
	b = append(b, []byte(Payload)...)
	fmt.Println(b)
	decode, err := Decode(b)
	fmt.Println(decode)
	if err != nil {
		t.Error(err)
	}
}

// 编码发送
func TestEncodeSubmit(t *testing.T) {
	var b []byte
	var ID = "12345678"
	var Payload = "hello"
	b = append(b, []byte(ID)...)
	b = append(b, []byte(Payload)...)
	s := &Submit{
		ID:      ID,
		Payload: []byte(Payload),
	}
	encode, err := Encode(s)
	fmt.Println(encode)
	if err != nil {
		t.Error(err)
	}
}

func TestDecodeSubmitAck(t *testing.T) {
	var b []byte
	var ID = "12345678"
	var Result = uint8(0)
	b = append(b, CommandSubmitAck)
	b = append(b, []byte(ID)...)
	b = append(b, byte(Result))
	fmt.Println(b)
	decode, err := Decode(b)
	fmt.Println(decode)
	if err != nil {
		t.Error(err)
	}
}

func TestEncodeSubmitAck(t *testing.T) {
	var b []byte
	var ID = "12345678"
	var Result = uint8(0)
	b = append(b, []byte(ID)...)
	b = append(b, byte(Result))
	s := &SubmitAck{
		ID:     ID,
		Result: Result,
	}
	encode, err := Encode(s)
	fmt.Println(encode)
	if err != nil {
		t.Error(err)
	}
}
