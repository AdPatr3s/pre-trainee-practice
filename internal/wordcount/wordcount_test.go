package wordcount_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/AdPatr3s/pre-trainee-practice/internal/wordcount"
)

func TestTokens(t *testing.T) {
	tests := []struct{
		name string
		in string
		want []string
	} {
		{name: "nil", in: "", want: nil},
		{name: "punctuation", in: "...", want: nil},
		{name: "normal", in: "Go, go, GO!", want: []string{"go", "go", "go"}},
		{name: "special", in: "naïve", want: []string{"naïve"}},
		{name: "plus", in: "C++", want: []string{"c"}},
		{name: "another", in: "it's", want: []string{"it", "s"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wordcount.Tokens(tt.in)
			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens(%s) = %v, want %v", tt.in, got, tt.want)
			}
			
		})
	}
}

func TestDropStop(t *testing.T) {
	tests := []struct{
		name string
		in []string
		stop map[string]struct{}
		want []string
	}{
		{name:"normal", in: []string{"a", "the", "go"}, stop: map[string]struct{}{"go": {}}, want: []string{"a", "the"}},
		{name:"no stop", in: []string{"a", "the", "go"}, stop: nil, want: []string{"a", "the", "go"}},
		{name:"empty tokens", in: []string{}, stop: map[string]struct{}{"go": {},}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		got := wordcount.DropStop(tt.in, tt.stop)
		
			if !slices.Equal(got, tt.want) {
				t.Errorf("DropStop(%v, %v) = %v, want %v", tt.in, tt.stop, got, tt.want)
			}
			
		})
	}
}

func TestTopN(t *testing.T) {
	tests := []struct{
		name string
		in map[string]int
		top int
		want []wordcount.Pair
	}{
		{name: "normal", in: map[string]int{"go": 4, "programming": 3, "clean": 1, "concise": 1}, top: 3, want: []wordcount.Pair{{"go", 4}, {"programming", 3}, {"clean", 1},}},
		{name: "normal", in: map[string]int{"go": 4, "programming": 3, "clean": 1, "concise": 1}, top: 1, want: []wordcount.Pair{{"go", 4},}},
		{name: "normal", in: map[string]int{"go": 4, "programming": 3, "clean": 1, "concise": 1}, top: -1, want: []wordcount.Pair{}},
		{name: "normal", in: nil, top: 1, want: []wordcount.Pair{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wordcount.TopN(tt.in, tt.top)
			if !slices.Equal(got, tt.want) {
				t.Errorf("TopN(%v, %v) = %v, want = %v", tt.in, tt.top, got, tt.want)
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want map[string]int
	}{
		{
			name: "normal",
			in:   []string{"a", "b", "a", "c", "a"},
			want: map[string]int{"a": 3, "b": 1, "c": 1},
		},
		{
			name: "empty tokens",
			in:   []string{},
			want: map[string]int{}, 
			
		},
		{
			name: "nil tokens",
			in:   nil,
			want: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wordcount.Count(tt.in)
			if !maps.Equal(got, tt.want) {
				t.Errorf("Count(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}