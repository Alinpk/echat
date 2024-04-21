package core

import (
	"net"
	"context"
	"sync"
	"serv/utils/log"
	"serv/proto"
	"serv/core/intf"
)



type User struct {
	// base info
	Conn net.Conn // net connection with client
	Certified bool // a flag to check if user has been login
	UserName, PassWord string // may password is not necessary
	
	// write tools
	// msg from server to client no need to reply for now.
	// -TODO who take the responsible for WriteBuffer
	// If the sender is responsible for the write buffer, a panic may be thrown when someone else writes to the write buffer.
	// solution: defer recover or use a rwlock(in most scenarios, we only use rlock, so this way could be effective,but i )
	WriteBuf chan ([]byte)  // stream wait to send to client, //sender has responsible for this resource
	isClosed bool
	mut sync.RWMutex

	// recv tools
	RcvBuf chan *proto.Message // receiver has responsible for this resource

	Ctx context.Context
	CancelFunc context.CancelFunc 
}

func (u *User) Quit() {
	u.CancelFunc() 
	u.Conn.Close() // make receiver quit
}

// net.Conn's Read & Write operation can only happen in this func
// in specific goroutine
func (u *User) Handle() {
	go u.Receiver()	// recv from conn, and push into rcvbuf
	go u.Sender() // send to conn

	for {
		select {
		case msg, ok := <- u.RcvBuf:
			if !ok {
				u.Quit()
			}
			u.Dispatch(msg)

		case _ = <- u.Ctx.Done():
			return
		}
	}
}

func (u *User) Write(buf []byte) {
	u.mut.RLock()
	defer u.mut.RUnlock()
	if !u.isClosed {
		u.WriteBuf <- buf
	}
}

func (u *User) Receiver() {
	defer close(u.RcvBuf)
	for {
		msg, err := proto.ReadMsg(u.Conn)
		if err != nil {
			log.L.Warn("Read msg failed", "detail", err.Error(), "addr", u.Conn.RemoteAddr().String())
			u.Quit()
			return
		}
		u.RcvBuf <- &msg
	}
}

func (u *User) Sender() {
	defer func() {
		u.mut.Lock()
		defer u.mut.Unlock()
		u.isClosed = true
		close(u.WriteBuf)
	}()
	for {
		select {
		case stream, _ := <- u.WriteBuf:
			for len(stream) != 0 {
				n, err := u.Conn.Write(stream)
				if err != nil {
					log.L.Warn("Read msg failed", "detail", err.Error(), "addr", u.Conn.RemoteAddr().String(), "user", u.UserName)
					u.Quit()
					return
				}
				stream = stream[n:]
			}
		case _ = <- u.Ctx.Done():
			return
		}
	}
}

func (u *User) Dispatch(msg *proto.Message) {
	// if type error
	defer func() {
		if r := recover(); r != nil {
			// TODO
			// send a service error msg
			// suppose this will never happened
			// log error
		}
	}()
	switch msg.Type {
	case REGISTER:
		u.Register(msg.Data.(*proto.RegisterMessage))
	case LOGIN:
		u.Register(msg.Data.(*proto.LoginMessage))
	case CONTROL:
		u.Control(msg.Data.(*proto.ControlMessage))
	case GROUP:
		u.Group(msg.Data.(*proto.GroupMessage))
	case PRIVATE:
		u.Group(msg.Data.(*proto.PrivateMessage))
	// server not support response
	default:
		log.L.Warn("get an unrecognized message")
		u.Write(BuildResponse(proto.INVALID, proto.BAD_REQUEST, ""))
	}
}

//----------------------service code-------------------------
func (u *User) Register(in *proto.RegisterMessage) {
	// only support before login
	if u.Certified {
		u.Write(BuildResponse(proto.REGISTER, proto.FORBIDDEN, "logout first"))
		return
	}

	ret, err := intf.RegisterUser(in.UserName, u.PassWord)
	if err != nil {
		log.L.Warn("register failed", "error", err.Error())
		u.Write(BuildResponse(proto.REGISTER, proto.INTERNAL_ERR, ""))
		return
	}
	if ret {
		log.L.Debug("register success", "user", in.UserName)
		u.Write(BuildResponse(proto.REGISTER, proto.OK, ""))
	} else {
		log.L.Info("register failed", "user", in.UserName)
		u.Write(BuildResponse(proto.REGISTER, proto.FORBIDDEN, "user existed"))
	}
}

func (u *User) Login(in *proto.LoginMessage) {
	// only support before login
	if u.Certified {
		u.Write(BuildResponse(proto.LOGIN, proto.FORBIDDEN, "logout first"))
		return
	}

	ret, err := intf.CheckUserAndPwd(in.UserName, in.PassWord)
	if err != nil {
		log.L.Warn("login failed", "error", err.Error())
		u.Write(BuildResponse(proto.LOGIN, proto.INTERNAL_ERR, ""))
		return
	}
	if ret {
		log.L.Debug("login success", "user", in.UserName)
		u.Write(BuildResponse(proto.LOGIN, proto.OK, ""))
		u.Certified = true
		u.UserName = in.UserName
		u.PassWord = in.PassWord
	} else {
		log.L.Info("login failed", "user", in.UserName)
		u.Write(BuildResponse(proto.LOGIN, proto.FORBIDDEN, "user existed"))
	}
}

func (u *User) Control(in *proto.ControlMessage) {
	// only support before login
	if !u.Certified {
		u.Write(BuildResponse(proto.CONTROL, proto.FORBIDDEN, "login first"))
		return
	}
	switch in.Type {
	case CREATE_ROOM:
		intf.C
	case  JOIN_ROOM:
	case QUIT_ROOM:
	}
}

func (u *User) Group(in *proto.GroupMessage) {
	// only support before login
	if !u.Certified {
		u.Write(BuildResponse(proto.GROUP, proto.FORBIDDEN, "login first"))
		return
	}

}

//TODO unrealize
func (u *User) Private(in *proto.PrivateMessage) {}


func BuildResponse(type proto.MsgType, code proto.RespCode, info string) []byte {
	ret, err := proto.EncodeMsg(ResponseMessage{
		Type : type,
		Code : code,
		Info : info,
	})
	if err != nil {
		log.L.Error("internal error", "error", err.Error())
		// empty msg
		return []byte{0, 0}
	}
	return ret
}