package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randRunes(n int, symbols []rune) []rune {
	runes := make([]rune, n)
	length := len(symbols)
	for i := range runes {
		runes[i] = symbols[rand.Intn(length)]
	}
	return runes
}

func GenerateDigits(length int) string {
	return string(randRunes(length, []rune("0123456789")))
}

func RandInt(min, max int) int {
	return rand.Intn(max-min) + min
}
