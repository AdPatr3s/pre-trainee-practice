package main

import (
	"fmt"

	"github.com/AdPatr3s/pre-trainee-practice/internal/fileproc"
)

func main() {


	

	fmt.Println(fileproc.Process([]string{"../../internal/fileproc/testdata/ok.txt", "../../internal/fileproc/testdata/missing.txt", "../../internal/fileproc/testdata/panic.txt", "../../internal/fileproc/testdata/empty.txt"}))
}	