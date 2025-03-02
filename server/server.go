package main

import (
	"ChatRoom/model"
	"fmt"
	"net"
)

var UserCenter = make(map[string]*model.User)
var Broadcast = make(chan string, 10)

func main() {
	fmt.Println("----聊天室服务器已启动----")
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		//有客户端访问服务器就Accept捕获
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		//进行处理
		go handler(conn)
	}

}

func broadcast() {
	for {
		msg := <-Broadcast

		for _, user := range UserCenter {
			user.MsgChan <- msg
		}
	}
}

func writeBackToClient(user *model.User) {
	for {
		msg := <-user.MsgChan
		_, _ = user.Conn.Write([]byte(msg))
	}
}

func handler(conn net.Conn) {
	name := conn.RemoteAddr().String()
	user := model.NewUser(name, conn)
	UserCenter[name] = user
	fmt.Println(name, "成功连接")
	info := fmt.Sprintf("%s进入了聊天室", name)
	Broadcast <- info
	go broadcast()
	go writeBackToClient(user)
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			// 检查错误类型，确定是否是连接关闭导致的
			if err.Error() == "EOF" {
				fmt.Println(name, "断开连接")
				info = fmt.Sprintf("%s离开了聊天室", name)
				Broadcast <- info
				return
			}
			delete(UserCenter, name) // 从用户中心删除已断开连接的用户
			return
		}
		info = fmt.Sprintf("%s:%s", name, string(buf[:n]))
		Broadcast <- info
	}
}
