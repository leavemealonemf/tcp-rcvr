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
	AVERAGE_IMEI_LEN = 10
	IMEI_LEN         = 15
)

const (
	network = ":8000"
)

// 003800360030003300380039003000350035003900320036003700330033

type Device struct {
	IMEI       string
	Socket     net.Conn
	Handshaked bool
}

var logger *log.Logger
var devices map[string]Device

func disconnectLog(locAddr, remAddr net.Addr) {
	logger.Printf("Disconnect\nLocal Addr: %v;\nRemote Addr: %v", locAddr, remAddr)
}

func handleConn(c net.Conn) {
	logger.Printf("Recieved new connection\nLocal Addr: %v\nRemote Addr: %v", c.LocalAddr(), c.RemoteAddr())

	var imei string
	var isImeiSaved bool = false

	defer c.Close()
	defer disconnectLog(c.LocalAddr(), c.RemoteAddr())

	reader := bufio.NewReader(c)

	for {
		msg, err := reader.ReadString('\n')

		if err != nil {
			logger.Println("Recieved msg err:\n", err.Error())
			if len(imei) > 10 {
				logger.Printf("Delete Device with imei: %v from map\n", imei)
				delete(devices, imei)
				isImeiSaved = false
			}
			break
		}

		logger.Println("Recieved msg:", msg)

		if !isImeiSaved {
			if len(msg) == HEXABLE_IMEI_LEN {
				dec, err := hex.DecodeHexStr(msg)
				fmt.Println("DECODED IMEI LEN", len(dec))
				if err != nil {
					logger.Println(err.Error())
					break
				}

				if len(dec) == IMEI_LEN {
					logger.Println("IMEI detected. Adding in map...")
					imei = dec
					devices[dec] = Device{
						IMEI:       dec,
						Socket:     c,
						Handshaked: true,
					}
					c.Write([]byte("01"))

					logger.Println("Decoded IMEI:", dec)
					logger.Println("IMEI Successfully saved")
					logger.Println("Saved DEC VAR IMEI:", dec)
					logger.Println("Saved IMEI:", imei)
					isImeiSaved = true
				} else {
					break
				}
			} else {
				break
			}
		} else {
			dec, err := hex.DecodeHexData(msg)
			if err != nil {
				logger.Println(err.Error())
			}
			logger.Printf("Decoded HEX Data From device with IMEI: %v\nDATA:%v\n", imei, dec)
			c.Write([]byte("Hello SIM800L from golang serve! =)"))
			time.Sleep(time.Second * 5)
		}
	}
}

func printDevices() {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("Saved devices:")
		fmt.Println("--------------")
		for k := range devices {
			fmt.Println(devices[k].IMEI)
		}
		fmt.Println("--------------")
	}
}

func main() {
	logger = log.Default()
	s, err := net.Listen("tcp", network)

	if err != nil {
		log.Fatalf("listen sock err: %v\n", err.Error())
	}

	devices = make(map[string]Device)

	for {
		conn, err := s.Accept()
		if err != nil {
			logger.Println("Received conn err:", err.Error())
		}
		go handleConn(conn)
		printDevices()
	}
}
