package rdb

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"errors"
	"reflect"
)

var id int
var pool *redis.Pool
func Get() redis.Conn {
	return pool.Get()
}

func LoadConfig(path string) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	var cfg redisCfg
	err = json.Unmarshal(text,&cfg)
	if err != nil {
		panic(err.Error())
	}
	id = cfg.DbId
	BuildConfig(cfg.Addr, cfg.Protocol, time.Duration(cfg.IdleTimeout) * time.Second, cfg.MaxIdle, cfg.MaxActive, cfg.DbId)
}

func BuildConfig(addr string, protocol string, idleTimeout time.Duration, maxIdle, maxActive, id int) {
	pool = &redis.Pool {
		MaxIdle: maxIdle,
		IdleTimeout: idleTimeout,
		MaxActive: maxActive,
		Dial: func() (redis.Conn, error) { return redis.Dial(protocol, addr, redis.DialDatabase(id)) },
	}

	poolConfigFix(pool)
	conn := Get()
	_, err := conn.Do("SELECT", id)
	if err != nil { panic(err.Error()) }
}


func poolConfigFix(pool *redis.Pool) {
	if pool == nil {
		return
	}
	if pool.MaxIdle < 0 { pool.MaxIdle = 50 }
	if pool.MaxActive < 0 { pool.MaxActive = 100 }
}

type redisCfg struct {
	Addr string		`json:"addr"`
	MaxIdle int		`json:"max_idle"`
	IdleTimeout int `json:"idle_timeout"`
	MaxActive int	`json:"max_active"`
	Protocol string	`json:"protocol"`
	DbId  int     	`json:"db_id"`
}

// if connect failed or something like this, return error
func HComapre(table, key, val string) (bool, error) {
	conn := Get()
	defer conn.Close()
	reply, err := ToStr(conn.Do("HGET", table, key))
	if err != nil { return false, err }
	return reply == val, nil
}

func HPutIfNotExisted(table, key, val string) (bool, error) {
	conn := Get()
	defer conn.Close()
	// if table not existed, then `exissted` == nil,
	// if table existed but key not existed, then `existed` == false
	expr := `
		local existed = redis.call('HGET', KEYS[1], KEYS[2])
		if existed == nil or existed == false
		then
			redis.call('HSET', KEYS[1], KEYS[2], KEYS[3])
			return true
		else
			return false
		end
	`
	luaScript := redis.NewScript(4, expr)
	reply, err := ToStr(luaScript.Do(conn, table, key, val, id))
	if err != nil { return false, err }
	// true->1, false->nil
	return reply == "1", nil
}

// !Please conduct thorough testing before use
func ToStr(reply interface{}, err error) (r string, e error) {
	if err != nil { e = err; return }
	if reply == nil { r = ""; return }
	switch v := reply.(type) {
	case string:
		r = v
	case int:
		r = strconv.Itoa(v)
	case int64:
		r = strconv.FormatInt(v, 10)
	case []byte:
		r = string(v)
	case bool:
		r = strconv.FormatBool(v)
	default:
		r = ""
		e = errors.New("unexpected type:" + reflect.TypeOf(reply).String())
	}
	return
}