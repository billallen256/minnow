package minnow

import (
	"fmt"
	"log"
)

type Dispatcher struct {
	WorkPath     Path
	IngestChan   chan IngestInfo
	ProcessorReg *ProcessorRegistry
	Logger       *log.Logger
}

func NewDispatcher(workPath Path, ingestChan chan IngestInfo, processorReg *ProcessorRegistry) (*Dispatcher, error) {
	if !workPath.Exists() {
		return nil, fmt.Errorf("Work path does not exist: %s", workPath)
	}

	logger := log.New(os.Stdout, "Dispatcher: ", 0)
	return &Dispatcher{workPath, ingestChan, processorReg, logger}
}

func (dispatcher *Dispatcher) Run() {
	for ingestInfo := range dispatcher.IngestChan {
		metadata, err := PropertiesFromFile(ingestInfo.MetadataPath)

		if err != nil {
			dispatcher.Logger.Print(err.Error())
			continue
		}

		matchingProcessors := dispatcher.ProcessorReg.MatchingProcessors(metadata)

		for _, processor := range matchingProcessors {
			inputPath := makeRandomPath(dispatcher.WorkPath).Resolve()
			outputPath := makeRandomPath(dispatcher.WorkPath).Resolve()
			//copy metadata
			err := Copy(metadataPath, inputPath)

			if err != nil {
				dispatcher.Logger.Print("Error copying metadata to work path: %s", err.Error())
				continue
			}

			//copy data
			err = Copy(dataPath, inputPath)
			if err != nil {
				dispatcher.Logger.Print("Error copying data to work path: %s", err.Error())
				continue
			}

			go processor.Run(inputPath, outputPath, ingestInfo.ProcessedBy, dispatcher.IngestChan)
		}
	}
}
