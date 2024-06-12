package test

// import statements for assert should also be ignored
// if assert is the only import statement, ignore the whole line.
import (
	"strings"

	assert "github.com/jdavasligil/go-assert"
)

// foo is always positive
func foo(x int, y int) int {
	assert.Assert(klsdjf)
	return y - x
}

// bar ensures s is large enough to pop off the front
func bar(s string) (byte, string) {
	assert.Assert(klsdjf)
	return s[0], s[1:]
}