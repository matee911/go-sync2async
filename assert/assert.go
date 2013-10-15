package assert

import (
	"bytes"
	"testing"
)

func Equal(t *testing.T, x interface{}, y interface{}) {
	// type cast i panic jak drugi typ sie nie zgadza
	switch x.(type) {
		case string: if x != y.(string) { t.Errorf("Expected: %v; Received: %v", y, x)}
		case byte: if !bytes.Equal(x.([]byte), y.([]byte)) { t.Errorf("Expected: %v; Received: %v", y, x) }
	}
}
