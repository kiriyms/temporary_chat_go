package utils

import "math/rand/v2"

func GenerateRoomCode(n int) string {
	digits := []string{
		"0",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
	}

	var res string
	for i := 0; i < n; i++ {
		res += digits[rand.IntN(len(digits))]
	}

	return res
}
