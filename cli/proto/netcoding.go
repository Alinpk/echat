package proto

import (
	"encoding/json"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func Decode(buf []byte, ptr interface{}) error {
	return json.Unmarshal(buf, ptr)
}

func Encode(rsp interface{}) ([]byte, error) {
	return json.Marshal(rsp)
}



func GetLen(data [2]byte) uint16 {
	return binary.BigEndian.Uint16(data[:])
}

func DecodeMsg(buf []byte) (msg Message, e error) {
	if len(buf) == 0 { e = errors.New("decode failed: len of buf = 0"); return }

	msg.Type = MsgType(buf[0])

	switch msg.Type {
	case REGISTER:
		msg.Data = &RegisterMessage{}
	case LOGIN:
		msg.Data = &LoginMessage{}
	case CONTROL:
		msg.Data = &ControlMessage{}
	case GROUP:
		msg.Data = &GroupMessage{}
	case PRIVATE:
		msg.Data = &PrivateMessage{}
	case RESPONSE:
		msg.Data = &ResponseMessage{}
	default:
		e = errors.New("unknown type")
		return
	}

	e = Decode(buf[1:], msg.Data)
	return
}

func EncodeMsg(i interface{}) (s []byte, e error) {
	var t MsgType
	switch i.(type) {
	case RegisterMessage:
		t = REGISTER
	case LoginMessage:
		t = LOGIN
	case ControlMessage:
		t = CONTROL
	case GroupMessage:
		t = GROUP
	case PrivateMessage:
		t = PRIVATE
	case ResponseMessage:
		t = RESPONSE
	default:
		e = errors.New("type is not supported")
		return
	}

	// len(2)+(type)
	prefix := make([]byte, MSG_LEN + MSG_TYPE_LEN)
	prefix[MSG_LEN] = byte(t)

	var msgData []byte
	if msgData, e = Encode(i); e != nil { return }

	// type + data
	realLen := len(msgData) + MSG_TYPE_LEN
	binary.BigEndian.PutUint16(prefix, uint16(realLen))

	s = append(prefix, msgData...)
	return
}

// maybe recv msg by this
func ReadMsg(r io.Reader) (msg Message, err error) {
	var length [2]byte
	var n int
	n, err = r.Read(length[:])
	fmt.Println("n:", n, "buf:", length)
	if err != nil { fmt.Println("err:", err.Error());return }

	l := int(GetLen(length))
	fmt.Println("expect get len:", l)
	buf := make([]byte, int(GetLen(length)))
	{
		var rdByte = 0
		for rdByte < len(buf) {
			n, e := r.Read(buf[rdByte:])
			if e != nil { err = e; fmt.Println("rdmsg:",err.Error()) }
			rdByte += n
			fmt.Println(rdByte)
		}
	}
	fmt.Println("Get full")
	msg, err = DecodeMsg(buf)
	return
}