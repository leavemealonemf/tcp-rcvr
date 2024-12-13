package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	liteDecoder "stcp/decoder/lite"
	"stcp/utils/hex"
	"strings"
	"time"
)

const (
	HEXABLE_IMEI_LEN = 62
	AVERAGE_IMEI_LEN = 10
	IMEI_LEN_ASCII   = 15
	IMEI_LEN_UNICODE = 30
)

const (
	network = ":8000"
)

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

				imeiLen := len(dec)

				fmt.Println("DECODED IMEI LEN", len(dec))
				if err != nil {
					logger.Println(err.Error())
					break
				}

				if imeiLen == IMEI_LEN_ASCII || imeiLen == IMEI_LEN_UNICODE {
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
			logger.Printf("Encoded HEX Data From device with IMEI: %v\nDATA:%v\n", imei, msg)
			clearPacket := hex.ClearInput(msg)
			packet, err := liteDecoder.DecodePacket(clearPacket)

			if err != nil {
				logger.Println(err.Error())
				continue
			}

			logger.Printf("Decoded JSON Data:\n\n%s\r\n", string(packet))

			c.Write([]byte("Hello SIM800L from golang serve! =)"))
			time.Sleep(time.Second * 5)
		}
	}
}

func main() {
	logger = log.Default()
	logger.SetPrefix("[SERVER] ")
	s, err := net.Listen("tcp", network)

	if err != nil {
		log.Fatalf("listen sock err: %v\n", err.Error())
	}

	logger.Printf("Waiting for accept connections...\nNetwork: %s\nAddr: %s\r\n", strings.ToUpper(s.Addr().Network()), s.Addr().String())

	devices = make(map[string]Device)

	for {
		conn, err := s.Accept()
		if err != nil {
			logger.Println("Received conn err:", err.Error())
		}
		go handleConn(conn)
	}
}
