package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	go listenFromSever(conn)
	inputReader := bufio.NewReader(os.Stdin)
	for {
		text, _ := inputReader.ReadString('\n')
		trimmedText := text[:len(text)-1] // 去除换行符
		_, _ = conn.Write([]byte(trimmedText))
	}
}

func listenFromSever(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("聊天室服务器已断开")
			return
		}
		fmt.Println(string(buf))
	}
}
