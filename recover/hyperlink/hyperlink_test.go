package hyperlink

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateLinks(t *testing.T) {
	tempOsName := osName
	stack := debug.Stack()
	links := CreateLinks(string(stack))
	assert.NotEqualf(t, "", links, "Expected links, got nothing.")

	defer func() {
		osName = tempOsName
	}()
	osName = "others"
	links = CreateLinks(string(stack))
	assert.NotEqualf(t, "", links, "Expected links, got nothing.")
}
