package core

import (
	"net"
	"sync/atomic"
	"context"
	"serv/proto"
	"serv/utils/log"
)

type Server struct {
	listener *net.Listener
	addr string
	active int32
}

func NewServer(addr string) *Server {
	listen, err := net.Listen("tcp", addr)
	log.L.Info("server starting", "addr", addr)
	if err != nil {
		log.L.Error("server init failed")
		panic("server init failed")
	}
	return  &Server{
		listener : &listen,
		addr : addr,
	}
}

func (s *Server) Start() {
	var consecutiveTimes int = 0
	for {
		conn, err := (*s.listener).Accept()
		if err != nil {
			consecutiveTimes++
			log.L.Warn("accept failed", "detail", err.Error(), "count", consecutiveTimes)
			if consecutiveTimes > 10 {
				log.L.Warn("quit")
				return
			}
		}
		consecutiveTimes = 0
		atomic.AddInt32(&s.active, 1)
		log.L.Info("new connection", "addr", conn.RemoteAddr().String(), "active", atomic.LoadInt32(&s.active))
		ctx, cancel := context.WithCancel(context.Background())
		var cli = &User{
			Conn : conn,
			Certified : false,
			WriteBuf : make(chan([]byte), 10),
			RcvBuf : make(chan(*proto.Message), 10),
			Ctx : ctx,
			CancelFunc : cancel,
			Groups : make(map[string](*Group)),
		}
		go func() {
			defer func() {
				atomic.AddInt32(&s.active, -1)
				log.L.Info("connection break", "addr", conn.RemoteAddr().String(), "active", atomic.LoadInt32(&s.active))
			}()
			cli.Handle()
		}()
	}
}