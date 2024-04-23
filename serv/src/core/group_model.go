package core

import (
	"sync"
	"strings"
	"reflect"
	"time"
	"unsafe"
	"serv/utils/file"
)

// key = groupName, value = *Group
var Groups sync.Map
const activeLimit = 500

type Group struct {
	Name string
	// < activeLimit
	UsersList []*User
	sync.RWMutex
	// TODO
	// type History interface { io.Writer }
	History *fileop.EFile
}

func NewGroup(name string) *Group {
	fs, err := fileop.OpenEFile("./cache/" + name)
	if err != nil {
		panic(err)
	}
	return &Group{
		Name : name,
		UsersList : make([]*User, 0),
		History : fs,
	}
}

func (g *Group) SpeakInGroup(user, msg string) {
	gmsg := strings.Join([]string {
		user,
		" ",
		time.Now().Format("2006-01-02 15:04:05"),
		"\n",
		msg,
		"\n",
	}, "")
	g.History.Write(GetSlice(gmsg))
	g.Write(BuildGroupMsg(g.Name, gmsg))
}

func (g *Group) Write(m []byte) {
	g.RLock()
	defer g.RUnlock()
	for _, ptr := range g.UsersList {
		ptr.Write(m)
	}
}

func (g *Group) QuitGroup(user *User) {
	g.DeleteUser(user)
	gmsg := strings.Join([]string {
		user.UserName,
		" quit ",
		time.Now().Format("2006-01-02 15:04:05"),
		"\n",
	}, "")
	g.History.Write(GetSlice(gmsg))
	g.Write(BuildGroupMsg(g.Name, gmsg))
}

func (g *Group) DeleteUser(user *User) {
	g.Lock()
	defer g.Unlock()
	l := len(g.UsersList)
	for i := 0; i < l; i++ {
		if g.UsersList[i] == user {
			// for gc
			g.UsersList[i] = nil
			g.UsersList[i], g.UsersList[l - 1] = g.UsersList[l - 1], g.UsersList[i]
			g.UsersList = g.UsersList[:l - 1]
			break
		}
	}
}

func (g *Group) AddUser(user *User) bool {
	if g.AddUserImpl(user) {
		gmsg := strings.Join([]string {
			user.UserName,
			" join ",
			time.Now().Format("2006-01-02 15:04:05"),
			"\n",
		}, "")
		g.History.Write(GetSlice(gmsg))
		g.Write(BuildGroupMsg(g.Name, gmsg))
		return true
	}
	return false
}

func (g *Group) AddUserImpl(user *User) bool {
	g.Lock()
	defer g.Unlock()
	if len(g.UsersList) >= activeLimit { return false }
	g.UsersList = append(g.UsersList, user)
	return true
}

func GetSlice(str string) []byte {
	// get ptr and len
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	strPtr := strHeader.Data
	strLen := strHeader.Len

	// build
	bytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: strPtr,
		Len:  strLen,
		Cap:  strLen,
	}))
	return bytes
}