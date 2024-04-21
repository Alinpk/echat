package main

import (
	"cli/proto"
	"fmt"
	"net"
	"sync"
)

func recv(conn net.Conn, wg sync.WaitGroup) {
	defer wg.Done()
	for {
		msg, err := proto.ReadMsg(conn)
		if err != nil {
			fmt.Println("error:", err.Error())
			break
		}
		fmt.Println("type:", msg.Type)
		fmt.Println("content:\n", msg.Data)
		fmt.Println("-------------------")
	}
}

func send(conn net.Conn, wg sync.WaitGroup) {
	defer wg.Done()
	for {
		var option string
		fmt.Scanf("%s\n", &option)

		switch option {
		case "register":
			var u, p string
			fmt.Scanf("%s\n%s\n", &u, &p)
			fmt.Println("user:", u)
			fmt.Println("password:", p)

			m := proto.RegisterMessage{
				UserName:u,
				PassWord:p,
			}

			msg, err := proto.EncodeMsg(m)
			if err != nil { panic(err.Error()) }
			var n int
			n, err = conn.Write(msg)
			if n != len(msg) || err != nil {
				fmt.Println("err:", err.Error(), n, "vs", len(msg))
			}


		case "login":
			var u, p string
			fmt.Scanf("%s\n%s\n", &u, &p)
			fmt.Println("user:", u)
			fmt.Println("password:", p)

			m := proto.LoginMessage{
				UserName:u,
				PassWord:p,
			}

			msg, err := proto.EncodeMsg(m)
			if err != nil { panic(err.Error()) }
			var n int
			n, err = conn.Write(msg)
			if n != len(msg) || err != nil {
				fmt.Println("err:", err.Error(), n, "vs", len(msg))
			}

		case "group":
			var t,c string
			fmt.Scanf("%s\n%s\n", &t, &c)
			m := proto.GroupMessage {
				GroupName : t,
				Content : c,
			}

			msg, err := proto.EncodeMsg(m)
			if err != nil { panic(err.Error()) }
			var n int
			n, err = conn.Write(msg)
			if n != len(msg) || err != nil {
				fmt.Println("err:", err.Error(), n, "vs", len(msg))
			}

		case "control":
			var op, target string
			fmt.Scanf("%s\n%s\n", &op, &target)
			m := proto.ControlMessage{
				Type : op,
				TargetName: target,
			}
			msg, err := proto.EncodeMsg(m)
			if err != nil { panic(err.Error()) }
			var n int
			n, err = conn.Write(msg)
			if n != len(msg) || err != nil {
				fmt.Println("err:", err.Error(), n, "vs", len(msg))
			}
			
		default:
			fmt.Println("invalid\n")
		}
	}
}


func main() {
	var wg sync.WaitGroup
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err.Error())
		return 
	}

	wg.Add(1)
	go recv(conn, wg)
	go send(conn, wg)
	wg.Wait()
}