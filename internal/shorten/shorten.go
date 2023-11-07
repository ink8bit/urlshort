package shorten

import (
	"math/rand"
	"strings"
)

// GenRandomStr creates randomized string from upper- and
// lowercased english letters and digits.
// The final length of the output string equals seven characters.
func GenRandomStr() string {
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
