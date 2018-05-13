package utils

import (
	"math/rand"
	"time"
)

// RemoveStrFromSlice attempts to remove a given string from a slice of
// strings.
func RemoveStrFromSlice(src []string, strMatch string) []string {
	for i, str := range src {
		if str == strMatch {
			src = append(src[:i], src[i+1:]...)
			break
		}
	}

	return src
}

// ShuffleStrings takes a source of strings and returns them in shuffled order.
func ShuffleStrings(src []string) []string {
	// Initialize a PRNG using the current time as a seed
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	dest := make([]string, len(src), len(src))
	perm := r.Perm(len(src))

	for i, v := range perm {
		dest[v] = src[i]
	}

	return dest
}
