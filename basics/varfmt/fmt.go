//go:build !solution

package varfmt

import (
	"fmt"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	ans := strings.Builder{}
	ptr := 0
	num := -2
	m := make(map[int]string, 0)
	for _, c := range format {
		// fmt.Println(i, c, num, ptr)
		if c == '{' {
			num = -1
		} else if num == -2 {
			ans.WriteRune(c)
		} else if '0' <= c && c <= '9' {
			if num == -1 {
				num = int(c - '0')
			} else {
				num *= 10
				num += int(c - '0')
			}
		} else if c == '}' {
			if num == -1 {
				val, exist := m[ptr]
				if !exist {
					m[ptr] = fmt.Sprint(args[ptr])
					val = m[ptr]
				}
				ans.WriteString(val)
			} else {
				val, exist := m[num]
				if !exist {
					m[num] = fmt.Sprint(args[num])
					val = m[num]
				}
				ans.WriteString(val)
			}
			num = -2
			ptr++
		}
	}
	return ans.String()
}
