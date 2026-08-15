package sliceutil_test

import (
	"testing"

	"flashquest/pkg/sliceutil"
)

func TestPtrSlice(t *testing.T) {
	items := []int{1, 2, 3}

	pointers := sliceutil.PtrSlice(items)

	if len(pointers) != len(items) {
		t.Fatalf("expected %d pointers, got %d", len(items), len(pointers))
	}

	for i, p := range pointers {
		if p != &items[i] {
			t.Fatalf("expected pointer %d to reference items[%d]", i, i)
		}
	}
}

func TestPtrSliceEmpty(t *testing.T) {
	pointers := sliceutil.PtrSlice([]string{})

	if pointers == nil {
		t.Fatal("expected non-nil slice")
	}
	if len(pointers) != 0 {
		t.Fatalf("expected empty slice, got %d elements", len(pointers))
	}
}
