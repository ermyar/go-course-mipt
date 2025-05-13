//go:build !solution

package hotelbusiness

import (
	"sort"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	m := make(map[int]int)
	for _, val := range guests {
		m[val.CheckInDate]++
		m[val.CheckOutDate]--
	}
	sum := 0
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	ans := make([]Load, 0, len(m))

	for _, key := range keys {
		value := m[key]
		sum += value
		if value == 0 {
			continue
		}
		ans = append(ans, Load{key, sum})
	}
	return ans
}
