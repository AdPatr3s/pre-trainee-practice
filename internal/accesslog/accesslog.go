package accesslog

import (
	"sort"
	"strconv"
	"strings"
)

type Line struct {
	IP     string
	Method string
	Path   string
	Status int
	Bytes  int
}

type Summary struct {
	Parsed    int
	Skipped   int
	UniqueIPs int
	ByStatus  map[int]int
	TopPaths  []PathCount // up to N, count desc, path asc on ties
}

type PathCount struct {
	Path  string
	Count int
}

// ParseLine parses one line.
// Returns ok=false for blank lines, comments, and malformed lines.
func ParseLine(s string) (Line, bool) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "#") || len(s) == 0 {
		return Line{}, false
	}
	results := strings.Split(s, " ")
	if len(results) != 5 {
		return Line{}, false
	}
	st, err1 := strconv.Atoi(results[3])
	b, err2 := strconv.Atoi(results[4])
	if !strings.HasPrefix(results[2], "/") || results[1] == "" || err1 != nil || err2 != nil {
		return Line{}, false
	}
	return Line{
		IP: results[0],
		Method: results[1],
		Path: results[2],
		Status: st,
		Bytes: b,
	}, true
}

// Summarize parses all lines (newline-separated).
// topN is the maximum length of TopPaths. If topN <= 0, TopPaths is empty.
func Summarize(text string, topN int) Summary {
	if text == "" {
		return Summary{ByStatus: map[int]int{}, TopPaths: []PathCount{}}
	}

	ips := make(map[string]struct{})
	byStatus := make(map[int]int)
	pathFreq := make(map[string]int)
	parsed, skipped := 0, 0

	for _, line := range strings.Split(text, "\n") {
		l, ok := ParseLine(line)
		if !ok {
			skipped++
			continue
		}
		parsed++
		ips[l.IP] = struct{}{}
		byStatus[l.Status]++
		pathFreq[l.Path]++
	}

	pairs := make([]PathCount, 0, len(pathFreq))
	for p, c := range pathFreq {
		pairs = append(pairs, PathCount{Path: p, Count: c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Path < pairs[j].Path
	})

	var top []PathCount
	if topN <= 0 {
		top = []PathCount{}
	} else if topN < len(pairs) {
		top = pairs[:topN]
	} else {
		top = pairs
	}

	return Summary{
		Parsed:    parsed,
		Skipped:   skipped,
		UniqueIPs: len(ips),
		ByStatus:  byStatus,
		TopPaths:  top,
	}
}