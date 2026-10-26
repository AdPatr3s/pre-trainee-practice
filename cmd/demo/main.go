package main

import (
	"errors"
	"fmt"

	// "os"

	// "github.com/AdPatr3s/pre-trainee-practice/internal/accesslog"
	// "github.com/AdPatr3s/pre-trainee-practice/internal/pipeline"
	// "github.com/AdPatr3s/pre-trainee-practice/internal/wordcount"
	"github.com/AdPatr3s/pre-trainee-practice/internal/booking"
)

func main() {
	// in := []int{1, 2, 3, 4, 5}
	// fmt.Println(pipeline.Double(in))
	// fmt.Println(pipeline.FilterMin(pipeline.Double(in), 5))
	// fmt.Println(pipeline.UniqueStable([]int{1, 2, 2, 3, 1, 4}))
	// fmt.Println(pipeline.UniqueStable([]int{7}))

	// s := "naïve C++ . Go is a programming language. Go is expressive, concise, clean, and efficient. Concurrency is a core idea in Go programming. Programming in Go is fun."
	// stop := map[string]struct{}{
	// 	"a": {}, "an": {}, "the": {}, "is": {}, "in": {},
	// 	"on": {}, "of": {}, "and": {}, "or": {}, "to": {}, "for": {},
	// }

	// tokens := wordcount.DropStop(wordcount.Tokens(s), stop)
	// freq := wordcount.Count(tokens)

	// fmt.Println(tokens)
	// fmt.Println(wordcount.TopN(freq, 5))

	// data, err := os.ReadFile("../../internal/accesslog/testdata/access.log")
	// if err != nil {
	// 	fmt.Println("Не удалось прочитать файл")
	// 	return
	// }
	// text := string(data)
	// fmt.Println(accesslog.Summarize(text, 3))

	s := booking.NewService()
	_ = s.AddRoom("orion", "Orion")

	a, _ := s.Book("orion", 9, 11)
	b, _ := s.Book("orion", 11, 12)

	if _, err := s.Book("orion", 10, 12); errors.Is(err, booking.ErrConflict) {
		fmt.Println("conflict detected")
	}

	_ = s.Cancel(a.ID)
	if _, err := s.Book("orion", 9, 10); err == nil {
		fmt.Println("slot freed after cancel")
	}

	list, _ := s.ListRoom("orion")
	fmt.Printf("reservations: %+v\n", list)
	_ = b
}	