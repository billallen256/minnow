package minnow

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
