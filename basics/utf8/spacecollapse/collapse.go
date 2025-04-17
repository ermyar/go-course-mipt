//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var b strings.Builder
	bytes := []byte(input)
	ok := 0
	for ptr := 0; ptr < len(bytes); {
		r, size := utf8.DecodeRune(bytes[ptr:])
		ptr += size
		if unicode.IsSpace(r) {
			if ok == 1 {
				continue
			} else {
				ok ^= 1
				b.WriteRune(0x20) // whitespace rune
			}
		} else {
			b.WriteRune(r)
			ok = 0
		}
	}
	return b.String()
}
