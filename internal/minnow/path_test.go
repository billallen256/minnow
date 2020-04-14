package minnow

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func randomString(length int) string {
	var builder strings.Builder
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		c := string(chars[rand.Intn(len(chars))])
		builder.WriteString(c)
	}

	return builder.String()
}

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

func TestIsFile(t *testing.T) {
	if !Path("/etc/passwd").IsFile() {
		t.Errorf("Path should be a file")
	}
}

func TestIsNotFile(t *testing.T) {
	if Path("/usr").IsFile() {
		t.Errorf("Path should not be a file")
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

func suffixTest(orig, newSuffix, target string, t *testing.T) {
	changed := Path(orig).WithSuffix(newSuffix)

	if changed != Path(target) {
		t.Errorf("WithSuffix failed: %s != %s", changed, target)
	}
}

func TestWithSuffix(t *testing.T) {
	suffixTest("foo.bar", "baz", "foo.baz", t)
	suffixTest("foo", "baz", "foo.baz", t)
	suffixTest("/foo/bar/baz.a", "baz", "/foo/bar/baz.baz", t)
	suffixTest("/foo/bar.a/baz.zip", "", "/foo/bar.a/baz", t)
}

func TestTouch(t *testing.T) {
	p := Path("/tmp/pathlib-" + randomString(20))
	p.Touch()

	if !p.Exists() {
		t.Errorf("Touch failed for path %s", p)
	}
}
