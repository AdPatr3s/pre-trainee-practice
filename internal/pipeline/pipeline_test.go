package pipeline_test

import (
	"slices"
	"testing"

	"github.com/AdPatr3s/pre-trainee-practice/internal/pipeline"
)

func TestDouble(t *testing.T) {
	tests := []struct{
		name string
		in []int
		want []int
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []int{}, want: []int{}},
		{name: "normal", in: []int{1, 2, 3}, want: []int{2, 4, 6}},
		{name:"negatives", in: []int{-1, 0}, want: []int{-2, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pipeline.Double(tt.in)
			if !slices.Equal(got, tt.want) {
				t.Errorf("Double(%v) = %v, want %v", tt.in, got, tt.want)
			}

			if tt.in == nil && got != nil {
				t.Errorf("Double(nil) returned non-nil: %v", got)
			}
		})
	}
}

func TestDouble_Aliasing(t *testing.T) {
	base := []int{1, 2, 3, 4}
	orig := slices.Clone(base)
	
	got := pipeline.Double(base[:2])
	if len(got) > 0 {
		got[0] = 99 
	}

	if !slices.Equal(base, orig) {
		t.Fatalf("Double aliased backing array: base=%v, want %v", base, orig)
	}
}

func TestFilterMin(t *testing.T) {
	tests := []struct{
		name string
		in []int
		min int
		want []int
	} {
		{name: "nil", 		in: nil, 			 min: 0,  want: nil},
		{name: "empty", 	in: []int{}, 		 min: 0,  want: []int{}},
		{name: "normal", 	in: []int{-1, 0, 3}, min: 0,  want: []int{0, 3}},
		{name: "zero", in: []int{1, 2, 3}, 		 min: 5,  want: []int{}}, 
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pipeline.FilterMin(tt.in, tt.min)
			if !slices.Equal(got, tt.want) {
				t.Errorf("FilterMin(%v, %v) = %v, want %v", tt.in, tt.min, got, tt.want)
			}
			
		})
	}
}

func TestUniqueStable(t *testing.T) {
	tests := []struct{
		name string
		in []int
		want []int
	} {
		{name: "nil", 		in: nil, 		     want: nil},
		{name: "empty", 	in: []int{}, 		 want: []int{}},
		{name: "test1", 	in: []int{1, 2, 2, 3, 1, 4}, want: []int{1, 2, 3, 4}},
		{name: "test2",      in: []int{2, 1, 2}, 		 want: []int{2, 1}}, 
		{name: "only one",      in: []int{7}, 		 want: []int{7}}, 
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pipeline.UniqueStable(tt.in)
			if !slices.Equal(got, tt.want) {
				t.Errorf("UniqueStable(%v) = %v, want %v", tt.in, got, tt.want)
			}
			
		})
	}
}