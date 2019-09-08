package hyperlink

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateLinks(t *testing.T) {
	stack := debug.Stack()
	links := CreateLinks(string(stack))
	assert.NotEqualf(t, "", links, "Expected links, got nothing.")
}
