//go:build !solution

package speller

import (
	"fmt"
	"strings"
)

// lookup-table for special numbers
var exclusiveNum = map[int]string{
	0:   "zero",
	1:   "one",
	2:   "two",
	3:   "three",
	4:   "four",
	5:   "five",
	6:   "six",
	7:   "seven",
	8:   "eight",
	9:   "nine",
	10:  "ten",
	11:  "eleven",
	12:  "twelve",
	13:  "thirteen",
	14:  "fourteen",
	15:  "fifteen",
	16:  "sixteen",
	17:  "seventeen",
	18:  "eighteen",
	19:  "nineteen",
	20:  "twenty",
	30:  "thirty",
	40:  "forty",
	50:  "fifty",
	60:  "sixty",
	70:  "seventy",
	80:  "eighty",
	90:  "ninety",
	100: "hundred",
}

func getSpelling(num []int) []string {
	tmp := make([]string, 0, 6)
	switch len(num) {
	case 3:
		if num[2] != 0 {
			tmp = append(tmp, exclusiveNum[num[2]], exclusiveNum[100])
		}
		fallthrough
	case 2:
		dig := num[1]*10 + num[0]
		if dig == 0 {
			return tmp
		}
		if num[0] == 0 {
			tmp = append(tmp, exclusiveNum[dig])
		} else if 21 <= dig {
			tmp = append(tmp, fmt.Sprintf("%s-%s", exclusiveNum[num[1]*10], exclusiveNum[num[0]]))
		} else {
			tmp = append(tmp, exclusiveNum[dig])
		}
	case 1:
		if num[0] == 0 {
			return tmp
		}
		tmp = append(tmp, exclusiveNum[num[0]])
	default:
		return tmp
	}
	return tmp
}

func Spell(n int64) string {
	if n == 0 {
		return exclusiveNum[0]
	}
	minus := false
	if n < 0 {
		minus = true
		n = -n
	}
	arr := make([]int, 0, 13)
	for n > 0 {
		arr = append(arr, int(n%10))
		n /= 10
	}
	ans := make([]string, 0)
	if minus {
		ans = append(ans, "minus")
	}
	for i, str := range []string{"billion", "million", "thousand", ""} {
		if tmp := getSpelling(arr[(9 - 3*i):(12 - 3*i)]); len(tmp) > 0 {
			ans = append(ans, tmp...)
			if i != 3 {
				ans = append(ans, str)
			}
		}
	}
	return strings.Join(ans, " ")
}
