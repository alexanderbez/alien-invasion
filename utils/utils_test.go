package utils

import (
	"reflect"
	"testing"
)

func TestRemoveStrFromSlice(t *testing.T) {
	testCases := []struct {
		s        []string
		e        []string
		strMatch string
	}{
		{
			s:        []string{"a", "b", "c", "d"},
			strMatch: "b",
			e:        []string{"a", "c", "d"},
		},
		{
			s:        []string{"a", "b", "c", "d"},
			strMatch: "z",
			e:        []string{"a", "b", "c", "d"},
		},
		{
			s:        []string{"a"},
			strMatch: "a",
			e:        []string{},
		},
	}

	for _, tc := range testCases {
		r := RemoveStrFromSlice(tc.s, tc.strMatch)

		if !reflect.DeepEqual(r, tc.e) {
			t.Errorf("incorrect result: expected: %v, got: %v", tc.e, r)
		}
	}
}
