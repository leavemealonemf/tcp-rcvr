package main

import (
	"bufio"
	"log"
	"net"
)

const (
	network = ":8080"
)

var logger *log.Logger
var socks []net.Conn

func handleConn(c net.Conn) {
	logger.Printf("Recieved new connection\nLocal Addr: %v\nRemote Addr: %v", c.LocalAddr(), c.RemoteAddr())

	defer c.Close()
	reader := bufio.NewReader(c)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			logger.Println("Recieved msg err:\n", err.Error())
			logger.Printf("Disconnect\nLocal Addr: %v;\nRemote Addr: %v", c.LocalAddr(), c.RemoteAddr())
			break
		}
		logger.Println("Recieved msg:", msg)
	}
}

func main() {
	logger = log.Default()
	s, err := net.Listen("tcp", network)

	if err != nil {
		log.Fatalf("listen sock err: %v\n", err.Error())
	}

	for {
		conn, err := s.Accept()
		if err != nil {
			logger.Println("Received conn err:", err.Error())
		}
		socks = append(socks, conn)
		go handleConn(conn)
	}
}
