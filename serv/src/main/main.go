package main

import (
	"serv/utils/log"
	"serv/utils/db"
	"serv/core"
	"fmt"
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
}