package hex

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func clearInput(input string) string {
	cleanInput := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1 // Убираем символ
		}
		return r
	}, input)

	return cleanInput
}

func DecodeHexData(h string) ([]int, error) {
	input := clearInput(h)

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

		num, err := strconv.ParseInt(dataPart, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования данных: %s", dataPart)
		}

		result = append(result, int(num))
	}

	return result, nil
}

func DecodeHexStr(hs string) (string, error) {
	input := clearInput(hs)
	dec, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования данных: %s", err.Error())
	}
	return string(dec), nil
}
