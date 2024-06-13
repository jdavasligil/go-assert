package test

// import statements for assert should also be ignored
// if assert is the only import statement, ignore the whole line.
import (
	"strings"
)

// foo is always positive
func foo(x int, y int) int {
	return y - x
}

// bar ensures s is large enough to pop off the front
func bar(s string) (byte, string) {
	return s[0], s[1:]
}
