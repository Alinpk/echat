package intf

// this file provide logical operation for model
import (
	"serv/utils/db"
	"serv/utils/log"
)

//------------------user--------------------
const (
	USER_TABLE string = "user_table"
	// maybe not use
	ONLINE_USER_TABLE string = "online_user_table"
)

func CheckUserAndPwd(user, pwd string) (bool, error) {
	return rdb.HComapre(USER_TABLE, user, pwd)
}


func RegisterUser(user, pwd string) (bool, error) {
	return HPutIfNotExisted(USER_TABLE, user, pwd)
}

//------------------grp--------------------
const (
	GROUP_TABLE string = "group_table"
)

// every group has a file to save it's own history
func genGrpFileName(groupName string) string {
	return "group_" + groupName
}

func NewGroup(groupName string) (bool, error) {
	return HPutIfNotExisted(GROUP_TABLE, groupName, genGrpFileName(groupName))
}