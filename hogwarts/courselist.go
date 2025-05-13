//go:build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) []string {
	visited := make(map[string]int)
	reverse := make(map[string][]string)
	for key, value := range prereqs {
		visited[key] = 1
		reverse[key] = append(reverse[key], "")
		for _, val := range value {
			visited[val] = 1
			reverse[val] = append(reverse[val], key)
		}
	}
	cnt := len(visited)
	stack := make([]string, cnt)
	lastQ := 0
	stats := make(map[string]int)
	for key := range reverse {
		stats[key] = len(prereqs[key])
		if stats[key] == 0 {
			stack[lastQ] = key
			lastQ++
		}
	}

	ans := make([]string, 0, cnt)

	for lastQ >= 0 {
		key := stack[lastQ]
		lastQ--
		ans = append(ans, key)
		visited[key] = 0
		cnt--
		for _, value := range reverse[key] {
			stats[value]--
			if stats[value] == 0 {
				lastQ++
				stack[lastQ] = value
			}
		}
	}
	if cnt > 0 {
		panic(-1)
	}
	return ans
}
