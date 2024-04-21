package proto

const (
	MSG_LEN = 2
	MSG_TYPE_LEN = 1
)

// definition of msg type 
type MsgType uint8
const (
	REGISTER MsgType = iota
	LOGIN
	CONTROL
	GROUP
	PRIVATE
	RESPONSE
	INVALID
)

// definition of response code
// refer to http resp code
type RespCode
const (
	BAD_REQUEST = 400
	UNAUTHORIZED = 401
	NOT_FOUND = 404
	FORBIDDEN = 403

	OK = 200

	INTERNAL_ERR = 500
	SERVIEC_UNAVAILABLE = 503
)

// // all pointer of msg should support decode and encode method
// type MsgT interface {
// 	Decode([]byte) error
// 	Encode() []byte
// 	Type() MsgType
// }

type Message struct {
	Type MsgType
	Data interface{}
}

/*  response msg
	@Type response for what type of msg
	@Code response code(Definitions are given above)
	@Error error msg(may empty str)
	@Data (if need)
*/
type ResponseMessage struct {
	Type MsgType
	Code RespCode
	Info string
}

// Register Msg
type RegisterMessage struct {
	UserName string
	PassWord string
}

// login msg
type LoginMessage struct {
	UserName string
	PassWord string
}

// chat in group
type GroupMessage struct {
	GroupName string
	Content string
}

// chat to person in private
type PrivateMessage struct {
	UserName string
	Content  string
}

// option suppose by control msg
const (
	CREATE_ROOM = "create_room"
	JOIN_ROOM = "join_room"
	QUIT_ROOM = "quit_room"
)
type ControlMessage struct {
	Type string
	TargetName string
}