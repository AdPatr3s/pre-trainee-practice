package main

import (
	"fmt"

	"github.com/AdPatr3s/pre-trainee-practice/internal/pipeline"
	"github.com/AdPatr3s/pre-trainee-practice/internal/wordcount"
)

func main() {
	in := []int{1, 2, 3, 4, 5}
	fmt.Println(pipeline.Double(in))
	fmt.Println(pipeline.FilterMin(pipeline.Double(in), 5))
	fmt.Println(pipeline.UniqueStable([]int{1, 2, 2, 3, 1, 4}))
	fmt.Println(pipeline.UniqueStable([]int{7}))

	s := "naïve C++ . Go is a programming language. Go is expressive, concise, clean, and efficient. Concurrency is a core idea in Go programming. Programming in Go is fun."
	stop := map[string]struct{}{
		"a": {}, "an": {}, "the": {}, "is": {}, "in": {},
		"on": {}, "of": {}, "and": {}, "or": {}, "to": {}, "for": {},
	}

	tokens := wordcount.DropStop(wordcount.Tokens(s), stop)
	freq := wordcount.Count(tokens)

	fmt.Println(tokens)
	fmt.Println(wordcount.TopN(freq, 5))

}	