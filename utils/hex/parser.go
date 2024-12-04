package hex

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func DecodeHexData(input string) ([]int, error) {
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
	dec, err := hex.DecodeString(hs)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования данных: %s", err.Error())
	}
	return string(dec), nil
}
