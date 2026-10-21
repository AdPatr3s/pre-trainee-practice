package pipeline

// Double returns a new slice where each element is n*2.
// If in == nil, return nil.
func Double(in []int) []int {
	if in == nil {
		return nil
	}

	result := make([]int, len(in))

	copy(result, in)
	for i := range in {
		result[i] = in[i] * 2
	}
	return  result
}

// FilterMin returns a new slice with elements >= min, in the same order.
// If in == nil, return nil.
func FilterMin(in []int, min int) []int {
	if in == nil {
		return nil
	}

	result := make([]int, 0)

	for _, v := range in {
		if v >= min {
			result = append(result, v)
		}
	}
	return result
}

// UniqueStable returns first occurrences only, original order kept.
// If in == nil, return nil.
func UniqueStable(in []int) []int {
	if in == nil {
		return nil
	}

	seen := make(map[int]struct{}, len(in))
	result := make([]int, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}