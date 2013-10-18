package dvs

import (
	"testing"
	"fmt"
	"github.com/matee911/go-sync2async/assert"
)

func v(f interface{}, err error) interface{} {
	return f
}

// Num

func TestEncodeNum(t *testing.T) {
	assert.Equal(t, EncodeNum(1, 1), "1")
	assert.Equal(t, EncodeNum(3, 2), "03")
	assert.Equal(t, EncodeNum(12, 2), "12")
	assert.Equal(t, EncodeNum(33, 3), "033")
	assert.Equal(t, EncodeNum(333, 3), "333")
	assert.Equal(t, EncodeNum(1234, 3), "1234")
}

func TestDecodeNum(t *testing.T) {
	assert.Equal(t, v(DecodeNum("1")), 1)
	assert.Equal(t, v(DecodeNum("03")), 3)
	assert.Equal(t, v(DecodeNum("12")), 12)
	assert.Equal(t, v(DecodeNum("033")), 33)
	assert.Equal(t, v(DecodeNum("333")), 333)
	assert.Equal(t, v(DecodeNum("0001")), 1)
}

func TestEncodeDecodeNum(t *testing.T) {
	assert.Equal(t, v(DecodeNum(EncodeNum(12, 2))), 12)
	assert.Equal(t, v(DecodeNum(EncodeNum(12, 3))), 12)
	assert.Equal(t, v(DecodeNum(EncodeNum(12, 4))), 12)
}

// Hex
func TestEncodeHex(t *testing.T) {
	assert.Equal(t, EncodeHex(1 ,1), []byte{1})
	assert.Equal(t, EncodeHex(65 ,2), []byte{0, 65})
}

func TestDecodeHex(t *testing.T) {
	assert.Equal(t, v(DecodeHex([]byte{0, 65})), 65)
}

func TestEncodeDecodeHex(t *testing.T) {
	assert.Equal(t, v(DecodeHex(EncodeHex(12, 1))), 12)
	assert.Equal(t, v(DecodeHex(EncodeHex(17, 2))), 17)
}

func TestRootHeader(t *testing.T) {
	assert.Equal(t, RootHeader(12, CmdTypeOther, 3, 75, 99), fmt.Sprint("000000012050003007599", CreationDate()))
}
