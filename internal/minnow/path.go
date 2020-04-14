package minnow

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Path string

func (p Path) Exists() bool {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return false
	}

	_, err = os.Stat(absPath)

	if err != nil {
		return false
	}

	return true
}

func (p Path) ReadBytes() ([]byte, error) {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadFile(absPath)

	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (p Path) IsDir() bool {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return false
	}

	stat, err := os.Stat(absPath)

	if err != nil {
		return false
	}

	mode := stat.Mode()

	if mode.IsDir() {
		return true
	}

	return false
}

func (p Path) Glob(pattern string) ([]Path, error) {
	if !p.IsDir() {
		return nil, fmt.Errorf("Glob only works on directories: %s", p)
	}

	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return nil, err
	}

	absPattern := filepath.Join(absPath, pattern)
	matches, err := filepath.Glob(absPattern)

	if err != nil {
		return nil, err
	}

	matchPaths := make([]Path, 0, len(matches))

	for _, match := range matches {
		matchPaths = append(matchPaths, Path(match))
	}

	return matchPaths, nil
}

func (p Path) Resolve() (Path, error) {
	absPath, err := filepath.Abs(string(p))

	if err != nil {
		return p, err
	}

	return Path(absPath), nil
}

func (p Path) WithSuffix(suffix string) Path {
	pStr := string(p)
	suffix = strings.TrimSpace(suffix)
	oldSuffix := filepath.Ext(pStr)

	if len(oldSuffix) == 0 && len(suffix) == 0 {
		return p
	}

	if len(oldSuffix) == 0 && len(suffix) > 0 {
		return Path(pStr + "." + suffix)
	}

	oldSuffixLen := len(oldSuffix)

	if len(suffix) > 0 {
		oldSuffixLen -= 1  // don't include the dot for removal
	}

	withoutOldSuffix := pStr[0:len(pStr)-oldSuffixLen]
	return Path(withoutOldSuffix + suffix)
}
