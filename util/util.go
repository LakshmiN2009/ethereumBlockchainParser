package util

import "fmt"

func HexToInt(hexStr string) int64 {
	var result int64
	_, err := fmt.Sscanf(hexStr, "0x%x", &result)
	if err != nil {
		fmt.Println("util.HexToInt Error converting hex to int: ", err)
		return 0
	}
	return result
}

func IntToHex(i int64) string {
	return fmt.Sprintf("0x%x", i)
}
