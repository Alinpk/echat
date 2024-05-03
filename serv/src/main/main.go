package main

import (
	"fmt"
	"serv/core"
	"serv/guard"
	rdb "serv/utils/db"
	"serv/utils/log"
	"time"
)

func main() {
	ResourceInit()
	server := core.NewServer(":9999")
	fmt.Println("server starting......")
	server.Start()
}

func ResourceInit() {
	path := "/home/huangzhujiang/gosdk/final/serv/config/"
	// db init
	rdb.LoadConfig(path + "db_cfg.json")
	// log init
	log.LoadCfg(path + "log_cfg.json")
	log.InitLog()
	guard.StartGuard(guard.GuardCfg{
		TidyUpInterval: time.Second * 3,
		GuardPath:      "./test/cache",
		FlowAddr:       "http://127.0.0.1:41555",
		ProcessDir:     "./test/backup",
	})
}
