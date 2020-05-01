package minnow

import (
	"log"
	"os"
	"time"
)

func Start(args []string) int {
	logger := log.New(os.Stdout, "Minnow: ", 0)

	if len(args) != 2 {
		logger.Print("Must specify a config file")
		return 1
	}

	config, err := ReadConfig(Path(args[1]))

	if err != nil {
		logger.Print(err.Error())
		return 1
	}

	dispatchChan := make(chan DispatchInfo, 1000)
	ingestDirChan := make(chan IngestDirInfo, 1000)
	defer close(ingestDirChan)
	defer close(dispatchChan)

	processorRegistry, err := NewProcessorRegistry(config.ProcessorDefinitionsPath)

	if err != nil {
		logger.Print(err.Error())
		return 1
	}

	dispatcher, err := NewDispatcher(config.WorkPath, dispatchChan, ingestDirChan, processorRegistry)

	if err != nil {
		logger.Print(err.Error())
		return 1
	}

	directoryIngester := NewDirectoryIngester(config.WorkPath, ingestDirChan, dispatchChan)

	go dispatcher.Run()
	go directoryIngester.Run()
	go processorRegistry.Run()

	for range time.Tick(config.IngestMinAge) {
		ingestDirChan <- IngestDirInfo{config.IngestPath, config.IngestMinAge, make([]ProcessorId, 0), false}
	}

	return 0
}
