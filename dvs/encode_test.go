package dvs

import (
	"testing"
	"fmt"
	"github.com/matee911/go-sync2async/assert"
)

func TestEnum(t *testing.T) {
	assert.Equal(t, Enum(1, 1), "1")
	assert.Equal(t, Enum(3, 2), "03")
	assert.Equal(t, Enum(12, 2), "12")
	assert.Equal(t, Enum(33, 3), "033")
	assert.Equal(t, Enum(333, 3), "333")
	assert.Equal(t, Enum(1234, 3), "1234")
}

func TestRootHeader(t *testing.T) {
	assert.Equal(t, RootHeader(12, CmdTypeOther, 3, 75, 99), fmt.Sprint("000000012050003007599", CreationDate()))
}
