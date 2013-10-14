package dvs

import (
	"testing"
	"fmt"
)

func assertS(t *testing.T, x string, y string) {
	if x != y {
		t.Errorf("Expected: %v; Received: %v", y, x)
	}
}

func TestEnum(t *testing.T) {
	assertS(t, Enum(1, 1), "1")
	assertS(t, Enum(3, 2), "03")
	assertS(t, Enum(12, 2), "12")
	assertS(t, Enum(33, 3), "033")
	assertS(t, Enum(333, 3), "333")
	assertS(t, Enum(1234, 3), "1234")
}

func TestRootHeader(t *testing.T) {
	assertS(t, RootHeader(12, CmdTypeOther, 3, 75, 99), fmt.Sprint("000000012050003007599", creationDate()))
}
