package minnow

import (
	"testing"
)

func TestExists(t *testing.T) {
	p := Path("/usr")

	if !p.Exists() {
		t.Errorf("Path should exist")
	}
}

func TestNotExists(t *testing.T) {
	p := Path("/foo/bar/baz/fjkdsalfjaklrejakfdsa")

	if p.Exists() {
		t.Errorf("Path should not exist")
	}
}

func TestReadBytes(t *testing.T) {
	p := Path("path.go")
	content, err := p.ReadBytes()

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(content) == 0 {
		t.Errorf("Received zero bytes")
	}
}
