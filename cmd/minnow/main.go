package main

import (
	"os"
	"github.com/gershwinlabs/minnow/internal/minnow"
)

func main() {
	os.Exit(minnow.Start(os.Args))
}
