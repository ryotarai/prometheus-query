package main

import (
	"fmt"
	"sort"
	"strings"
)

func stringMapToString(m map[string]string, delimiter string) string {
	s := make([][]string, len(m))
	i := 0
	for k, v := range m {
		s[i] = []string{k, v}
		i++
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i][0] < s[j][0]
	})

	ss := make([]string, len(s))
	for i, v := range s {
		ss[i] = fmt.Sprintf("%s:%s", v[0], v[1])
	}

	return strings.Join(ss, delimiter)
}
