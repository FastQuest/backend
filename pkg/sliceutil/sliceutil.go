package sliceutil

// PtrSlice returns a slice of pointers to the elements of items.
func PtrSlice[T any](items []T) []*T {
	out := make([]*T, len(items))
	for i := range items {
		out[i] = &items[i]
	}
	return out
}
