package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const limit = 10

var re = regexp.MustCompile(`[\s\t\n\r\.,:!\?]`)

func Top10(s string) []string {
	m := map[string]int{}
	for _, w := range re.Split(s, -1) {
		switch w {
		case "", "-":
			continue
		default:
			m[strings.ToLower(w)]++
		}
	}

	l := len(m)
	res := make([]string, 0, l)
	for w := range m {
		res = append(res, w)
	}

	sort.Slice(res, func(i, j int) bool {
		curr, next := res[i], res[j]
		if m[curr] == m[next] {
			return curr < next
		}
		return m[curr] > m[next]
	})

	if l > limit {
		return res[:limit]
	}
	return res
}
