//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	bytes := []byte(input)
	var b strings.Builder
	m := make(map[rune]int, 0)
	for i := len(bytes); i > 0; {
		r, size := rune(0), -1
		for ptr, tmp := i-1, rune(bytes[i-1]); ptr >= 0; ptr, tmp = ptr-1, (tmp<<8)+rune(bytes[ptr]) {
			if val, exist := m[tmp]; exist {
				r, size = tmp, val
				break
			}
		}
		if size == -1 {
			r, size = utf8.DecodeLastRune(bytes[:i])
			m[r] = size
		}
		if r == utf8.RuneError {
			b.WriteRune(utf8.RuneError)
		} else {
			b.Write(bytes[i-size : i])
		}
		i -= size
	}
	return b.String()
}
