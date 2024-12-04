package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"stcp/utils/hex"
	"time"
)

const (
	HEXABLE_IMEI_LEN = 62
)

const (
	network = ":8000"
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

		// DECODE MESSAGES
		fmt.Println("Msg len:", len(msg))
		if len(msg) == HEXABLE_IMEI_LEN {
			dec, err := hex.DecodeHexStr(msg)
			if err != nil {
				logger.Println(err.Error())
			}
			logger.Println(dec)
		} else {
			dec, err := hex.DecodeHexData(msg)
			if err != nil {
				logger.Println(err.Error())
			}
			logger.Println(dec)
		}

		c.Write([]byte("Hello SIM800L from golang serve! =)"))
		time.Sleep(time.Second * 5)
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
