# Wordcount Benchmarks

Environment: Windows, AMD Ryzen 5 3550H with Radeon Vega Mobile Gfx, Go 1.2x

## Results

BenchmarkJoinPlus-8                  536           2200049 ns/op         2689255 B/op       1000 allocs/op
BenchmarkJoinBuilder-8             42673             31615 ns/op           17912 B/op         13 allocs/op
BenchmarkTokens-8                    578           2355501 ns/op          467385 B/op       8514 allocs/op
BenchmarkCount-8                    2922            406491 ns/op             456 B/op          3 allocs/op

## Analysis

The `joinBuilder` is significantly faster and more memory-efficient than `joinPlus` (31,615 ns/op vs 2,200,049 ns/op). 
In Go, strings are immutable. Using `+=` in a loop creates a new string on every iteration and copies the entire previous string into the new one. 
This results in O(N^2) time complexity and allocates a huge amount of memory (1000 allocations for 1000 parts). 
`strings.Builder` uses a mutable byte slice under the hood, which grows amortized O(1), resulting in just a few allocations (13 allocs/op) and 70x faster execution. 

For `Tokens`, most of the 8514 allocations come from calling `strings.ToLower()` and `b.String()` for every single word, which creates a new string in memory. While pre-allocating the `strings.Builder` buffer with `b.Grow(len(s))` helps reduce buffer reallocations, the fundamental cost of returning a slice of new lowercased strings remains.