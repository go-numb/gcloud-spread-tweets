package models

import (
	"math"
	"strings"
)

// StrUpsideDown は文字列を伏せ字にします。
func StrUpsideDown(s string) string {
	l := len([]rune(s))
	if l <= 2 {
		return s
	}

	suffN := int(math.Floor(float64(l)/float64(2))) + 1
	suff := []rune(s)[suffN:]

	down := make([]string, suffN)
	for i := 0; i < suffN; i++ {
		down[i] = "*"
	}

	return strings.ReplaceAll(s, string(suff), strings.Join(down, ""))
}
