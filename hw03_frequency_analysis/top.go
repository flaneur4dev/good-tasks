package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type word struct {
	w string
	c int
}

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

	words := []word{}
	for k, v := range m {
		words = append(words, word{k, v})
	}

	sort.Slice(words, func(i, j int) bool {
		if words[i].c == words[j].c {
			return words[i].w < words[j].w
		}
		return words[i].c > words[j].c
	})

	limit := len(words)
	if limit > 10 {
		limit = 10
	}

	res := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		res = append(res, words[i].w)
	}

	return res
}
