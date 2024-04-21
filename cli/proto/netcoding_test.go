package proto

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func TestSerialize(test *testing.T) {
	login := LoginMessage{
		UserName : "alan",
		PassWord : "1123456",
	}

	buf, err := EncodeMsg(login)
	assert.Equal(test, err, nil)

	msg, e := DecodeMsg(buf[2:])
	// no error
	assert.Equal(test, e, nil)
	// decode type should same with encode type
	assert.Equal(test, msg.Type, LOGIN)
	// same type
	assert.Equal(test, reflect.TypeOf(msg.Data), reflect.TypeOf(&login))
	// same value
	decodeRes := *msg.Data.(*LoginMessage)
	assert.Equal(test, decodeRes, login)
}