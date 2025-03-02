package model

import "net"

type User struct {
	Username string `json:"username"`
	MsgChan  chan string
	Conn     net.Conn
}

func NewUser(username string, conn net.Conn) *User {
	return &User{
		Username: username,
		MsgChan:  make(chan string),
		Conn:     conn,
	}
}
