package dvs

import (
	"testing"
	"github.com/matee911/go-sync2async/assert"
)

func TestNoCommandBody(t *testing.T) {
	assert.Equal(t, NoCommand().Body, "1002")
}

func TestNoCommandType(t *testing.T) {
	assert.Equal(t, NoCommand().Type, CmdTypeOther)
}

func TestPushVodCommandType(t *testing.T) {
	assert.Equal(t, PushVodCommand().Type, CmdTypeVod)
}
