package shorten

import (
	"math/rand"
	"strings"
)

// GenStr creates randomized string from letters and digits.
func GenStr() string {
	defaultSize := 7
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	chars := []rune(lower + upper + digits)
	var sb strings.Builder
	for i := 0; i < defaultSize; i++ {
		sb.WriteRune(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}
