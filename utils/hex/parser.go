package hex

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func ClearInput(input string) string {
	cleanInput := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, input)

	return cleanInput
}

func DecodeHexData(h string) ([]int, error) {
	input := ClearInput(h)

	if len(input)%4 != 0 {
		return nil, fmt.Errorf("некорректная длина строки, должно быть кратно 4")
	}

	var result []int

	var iteration uint8 = 0

	for i := 0; i < len(input); i += 4 {
		iteration += 1

		if (iteration%2) == 0 && input[i:i+4] == "0000" {
			continue
		}

		dataPart := input[i : i+4]

		num, err := strconv.ParseInt(dataPart, 16, 16)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования данных: %s", dataPart)
		}

		result = append(result, int(num))
	}

	return result, nil
}

func DecodeHexStr(hs string) (string, error) {
	input := ClearInput(hs)
	dec, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования данных: %s", err.Error())
	}
	return string(dec), nil
}

func HexToFloat(hexValue string) (float32, error) {
	bytes, err := hex.DecodeString(hexValue)
	if err != nil {
		return 0, fmt.Errorf("ошибка декодирования HEX: %v", err)
	}

	if len(bytes) != 4 {
		return 0, fmt.Errorf("некорректная длина данных, ожидалось 4 байта, получено: %d", len(bytes))
	}

	// Преобразуем байты в uint32
	bits := uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])

	// Преобразуем uint32 в float32 (IEEE 754)
	floatValue := math.Float32frombits(bits)
	return floatValue, nil
}
