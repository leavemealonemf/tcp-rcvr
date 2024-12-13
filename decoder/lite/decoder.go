package liteDecoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type avlDataValue struct {
	ParamName string
	Value     uint32
}

type finalPacket struct {
	Timestamp  uint64 `json:"_ts"`
	Online     bool   `json:"online"`
	Lat        uint32 `json:"lat"`
	Lon        uint32 `json:"lon"`
	Speed      uint32 `json:"speed"`
	Altitude   uint32 `json:"altitude"`
	Angle      uint32 `json:"angle"`
	Satellites uint32 `json:"sat-count"`
	GSMSignal  uint32 `json:"signal"`
}

var packetMap map[string]*avlDataValue

func getTimestamp(timestamp string) time.Time {
	var timestampMs int64
	fmt.Sscanf(timestamp, "%x", &timestampMs)

	timestampSeconds := timestampMs / 1000
	millisecondsPart := timestampMs % 1000

	return time.Unix(timestampSeconds, millisecondsPart*1e6).UTC()
}

func hexToDec(chunk string) (int64, error) {
	num, err := strconv.ParseInt(chunk, 16, 16)
	if err != nil {
		return 0, fmt.Errorf("hex to decimal failed convert: %s", chunk)
	}

	return num, nil
}

func saveParam(id, value string) error {
	decVal, _ := hexToDec(value)
	val := packetMap[id]

	if val != nil {
		packetMap[id].Value = uint32(decVal)
		return nil
	}

	return fmt.Errorf("add param failed")
}

func DecodePacket(packet string) ([]byte, error) {
	if len(packet) == 0 {
		return nil, errors.New("recived empty packet")
	}

	atlTimestamp := getTimestamp(packet[20:36])

	atlLat, _ := hexToDec(packet[38:46])
	atlLon, _ := hexToDec(packet[46:54])
	atlAltitude, _ := hexToDec(packet[54:58])
	atlAngle, _ := hexToDec(packet[58:62])
	atlSatellites, _ := hexToDec(packet[62:64])
	atlSpeed, _ := hexToDec(packet[64:68])

	n1BytesProps, _ := hexToDec(packet[76:80])

	iteration := 80

	if n1BytesProps > 0 {
		offset := 0

		for offset < int(n1BytesProps) {
			id := packet[iteration : iteration+4]
			param := packet[iteration+4 : iteration+6]
			saveParam(id, param)
			offset += 1
			iteration += 6
		}
	}

	n2BytesProps, _ := hexToDec(packet[iteration : iteration+4])
	iteration += 4

	if n2BytesProps > 0 {
		offset := 0

		for offset < int(n2BytesProps) {
			id := packet[iteration : iteration+4]
			param := packet[iteration+4 : iteration+8]
			saveParam(id, param)
			offset += 1
			iteration += 8
		}
	}

	n4BytesProps, _ := hexToDec(packet[iteration : iteration+4])
	iteration += 4

	if n4BytesProps > 0 {
		offset := 0

		for offset < int(n4BytesProps) {
			id := packet[iteration : iteration+4]
			param := packet[iteration+4 : iteration+12]
			saveParam(id, param)
			offset += 1
			iteration += 12
		}
	}

	n8BytesProps, _ := hexToDec(packet[iteration : iteration+4])
	iteration += 4

	if n8BytesProps > 0 {
		offset := 0

		for offset < int(n8BytesProps) {
			id := packet[iteration : iteration+4]
			param := packet[iteration+4 : iteration+20]
			saveParam(id, param)
			offset += 1
			iteration += 20
		}
	}

	finalPacket := &finalPacket{
		Timestamp:  uint64(atlTimestamp.UnixMilli()),
		Online:     true,
		Lat:        uint32(atlLat),
		Lon:        uint32(atlLon),
		Speed:      uint32(atlSpeed),
		Altitude:   uint32(atlAltitude),
		Angle:      uint32(atlAngle),
		Satellites: uint32(atlSatellites),
		GSMSignal:  packetMap["0034"].Value,
	}

	finalPktJson, err := json.Marshal(finalPacket)

	if err != nil {
		return nil, errors.New("failed convert to JSON")
	}

	return finalPktJson, nil
}

func init() {
	packetMap = make(map[string]*avlDataValue)

	packetMap["0034"] = &avlDataValue{
		ParamName: "signal",
	}
}
