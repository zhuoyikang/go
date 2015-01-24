// Package stringutil contains utility functions for working with strings.
package stringutil

import (
	"fmt"
)

// Reverse returns its argument string reversed rune-wise left to right.
func Good(s string) {
	fmt.Println(Reverse(s))
}

//
func init() {
	// Good("good init")
	// fmt.Printf("stringutil good init\n")
}
