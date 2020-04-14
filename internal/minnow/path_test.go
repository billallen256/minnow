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

func TestIsDir(t *testing.T) {
	if !Path("/usr").IsDir() {
		t.Errorf("Path should be a directory")
	}
}

func TestIsNotDir(t *testing.T) {
	if Path("/etc/passwd").IsDir() {
		t.Errorf("Path should not be a directory")
	}
}

func TestGlob(t *testing.T) {
	paths, err := Path("/etc").Glob("*.conf")

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(paths) == 0 {
		t.Errorf("/etc/*.conf returned no results")
	}

	for _, path := range paths {
		absPath, err := path.Resolve()

		if err != nil {
			t.Errorf(err.Error())
		}

		if path != absPath {
			t.Errorf("Glob should only return absolute paths")
		}
	}
}
