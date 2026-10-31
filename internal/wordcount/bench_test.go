package wordcount

import (
	"strings"
	"testing"
)

func joinPlus(parts []string) string {
	var s string
	for _, p := range parts {
		s += p
	}
	return s
}

func joinBuilder(parts []string) string {
	var b strings.Builder 
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

func TestJoin(t *testing.T) {
	tests := []struct{
		name string
		in []string
		want string
	} {{name: "normal", in: []string{"aaa", "vvv"}, want: "aaavvv"},
	   {name: "normal2", in: []string{"cvb, asd", "wsd, lasd. we"}, want: "cvb, asdwsd, lasd. we"},
	   {name: "empty", in: []string{}, want: ""}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := joinPlus(tt.in)
			got2 := joinBuilder(tt.in)
			if got1 != got2 {
				t.Errorf("joinPlus(in) = %s != goinBuilder(in) = %s", got1, got2)
			}
		})
	}
}

var sink string

func BenchmarkJoinPlus(b *testing.B) {
	parts := make([]string, 1000)
	for i := range 1000 {
		parts[i] = "words"
	}

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		sink = joinPlus(parts)
	}
}

func BenchmarkJoinBuilder(b *testing.B) {
	parts := make([]string, 1000)
	for i := range 1000 {
		parts[i] = "words"
	}

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		sink = joinBuilder(parts)
	}
}


func BenchmarkTokens(b *testing.B) {
	paragraph := "Go is a programming language. Go is expressive, concise, clean, and efficient. "
	s := strings.Repeat(paragraph, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Tokens(s)
	}
}

func BenchmarkCount(b *testing.B) {
	paragraph := "Go is a programming language. Go is expressive, concise, clean, and efficient. "
	s := strings.Repeat(paragraph, 500)

	tokens := Tokens(s)

	b.ResetTimer()

	for i:= 0; i <= b.N; i++ {
		_ = Count(tokens)
	}
}