package model

import "github.com/bwmarrin/snowflake"

//User means a people who use the system.
type User struct {
	UID     int64
	Account string
	Passwd  string
}

//NewUser returns a User who UID generate by snowlake Algorithm
func NewUser(account string, passwd string) *User {
	node, _ := snowflake.NewNode(1)
	return &User{
		UID:     node.Generate().Int64(),
		Account: account,
		Passwd:  passwd,
	}
}
