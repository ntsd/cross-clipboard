package clipboard

import (
	"reflect"
	"testing"
)

func TestLimitAppend(t *testing.T) {
	tests := []struct {
		limit  int
		arr    [][]byte
		new    []byte
		result [][]byte
	}{
		{
			limit:  4,
			arr:    [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")},
			new:    []byte("e"),
			result: [][]byte{[]byte("b"), []byte("c"), []byte("d"), []byte("e")},
		},
		{
			limit:  3,
			arr:    [][]byte{[]byte("a"), []byte("b")},
			new:    []byte("c"),
			result: [][]byte{[]byte("a"), []byte("b"), []byte("c")},
		},
	}

	for _, test := range tests {
		got := limitAppend(test.limit, test.arr, test.new)

		if !reflect.DeepEqual(got, test.result) {
			t.Errorf("got %q, wanted %q", got, test.result)
		}
	}
}
