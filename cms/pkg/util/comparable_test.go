package util

import (
	"fmt"
	"testing"
)

func TestIsEqualityMap(t *testing.T) {
	fmt.Println(IsEqualityMap(
		map[string]string{"a": "a", "b": "b"},
		map[string]string{"a": "a", "b": "b"},
	))

	fmt.Println(IsEqualityMap(
		map[string]string{"a": "a", "b": "b"},
		map[string]string{"a": "aa", "b": "bb"},
	))

	fmt.Println(IsEqualityMap(
		map[string]string{"a": "a"},
		map[string]string{"a": "a", "b": "b"},
	))

	fmt.Println(IsEqualityMap(
		map[string]string{"a": "a", "c": "c"},
		map[string]string{"a": "a", "b": "b"},
	))
}
