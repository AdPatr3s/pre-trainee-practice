package main

import (
	"fmt"

	"github.com/AdPatr3s/pre-trainee-practice/internal/pipeline"
)

func main() {
	in := []int{1, 2, 3, 4, 5}
	fmt.Println(pipeline.Double(in))
	fmt.Println(pipeline.FilterMin(pipeline.Double(in), 5))
	fmt.Println(pipeline.UniqueStable([]int{1, 2, 2, 3, 1, 4}))
	fmt.Println(pipeline.UniqueStable([]int{7}))
}	