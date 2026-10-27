package main

import (
	"errors"
	"fmt"

	"github.com/AdPatr3s/pre-trainee-practice/internal/booking"
	"github.com/AdPatr3s/pre-trainee-practice/internal/fileproc"
)

func main() {


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

	fmt.Println(fileproc.Process([]string{"../../internal/fileproc/testdata/ok.txt", "../../internal/fileproc/testdata/missing.txt", "../../internal/fileproc/testdata/panic.txt", "../../internal/fileproc/testdata/empty.txt"}))
}	