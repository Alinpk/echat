package proto

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"fmt"
)

func Test1(test *testing.T) {
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

func Test2(test *testing.T) {
	ctl := ControlMessage{
		Type : "join_group",
		TargetName : "mygrp",
	}

	buf, err := EncodeMsg(ctl)
	assert.Equal(test, err, nil)

	msg, e := DecodeMsg(buf[2:])
	// no error
	assert.Equal(test, e, nil)
	// decode type should same with encode type
	assert.Equal(test, msg.Type, CONTROL)
	// same type
	assert.Equal(test, reflect.TypeOf(msg.Data), reflect.TypeOf(&ctl))
	// same value
	decodeRes := *msg.Data.(*ControlMessage)
	assert.Equal(test, decodeRes, ctl)
}