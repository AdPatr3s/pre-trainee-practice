package wordcount

import (
	"sort"
	"strings"
	"unicode"
)

type Pair struct {
	Word  string
	Count int
}

// Tokens splits s into words using the tokenizer rules below.
func Tokens(s string) []string {
	var result []string
	var b strings.Builder

	wr := func() {
		if b.Len() > 0 {
			result = append(result, strings.ToLower(b.String()))
			b.Reset()
		} 
	}
	
	for _, i := range s {
		if unicode.IsLetter(i) || unicode.IsDigit(i) {
			b.WriteRune(i)
		} else {
			wr()
		}
	}
	wr()
	return result
}

// Count returns frequency of each token.
func Count(tokens []string) map[string]int {
	result := make(map[string]int)
	for _, t := range tokens {
		result[t]++
	}
	return result
}

// TopN returns up to n pairs, highest count first.
// Ties: alphabetical ascending by Word.
// If n <= 0, return an empty slice.
func TopN(freq map[string]int, n int) []Pair {
	if n <= 0 {
		return []Pair{}
	}

	result := make([]Pair, 0, n)
	for k, v := range freq {
		p := Pair{k, v}
		result = append(result, p)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Word < result[j].Word
	})
	if len(result) > n {
		result = result[:n]
	}
		return result
}

// DropStop removes tokens that appear in stop (lowercase keys).
func DropStop(tokens []string, stop map[string]struct{}) []string {
	result := make([]string, 0, len(tokens))

	for _, v := range tokens {
		if _, ok := stop[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}