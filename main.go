package main

import (
	"fmt"
	"net"
)

func main() {
	server := createServer()
	go listenForCommands(server)

	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		fmt.Println("failed to start the server: ", err.Error())
	}
	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("failed to connect to client: ", err.Error())
		}

		fmt.Println("client connected: ", conn.RemoteAddr().String())
		go server.createClient(conn)
	}
}
