package minnow

import (
	"fmt"
)

func start(args []string) int {
	if len(args) != 1 {
		fmt.Println("Must specify a config file")
		return 1
	}

	config, err := ReadConfig(args[1])

	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	ingestChan := make(chan IngestInfo, 1000)
}
