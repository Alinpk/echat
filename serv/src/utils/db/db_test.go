package rdb


import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gomodule/redigo/redis"
)


func TestConfig(test *testing.T) {
	PATH := "./db_cfg.json"
	LoadConfig(PATH)
	conn := Get()
	defer func() {
		redis.String(conn.Do("FLUSHDB"))
		conn.Close()
	}()
	{
		assert.Equal(test, id, 15)

		reply, err := redis.String(conn.Do("SELECT", id))
		assert.Equal(test, err, nil)
		assert.Equal(test, reply, "OK")
	}
	{
		reply, err := redis.String(conn.Do("FLUSHDB"))
		assert.Equal(test, err, nil)
		assert.Equal(test, reply, "OK")

		r, e := HPutIfNotExisted("test", "caseNum", "adsd")
		assert.Equal(test, e, nil)
		assert.Equal(test, r, true)

		r, e = HComapre("test", "caseNum", "adsd")
		assert.Equal(test, e, nil)
		assert.Equal(test, r, true)

		r, e = HPutIfNotExisted("test", "caseNum", "adsd")
		assert.Equal(test, e, nil)
		assert.Equal(test, r, false)
	}
}