package minnow

import (
	"fmt"
	"log"
	"os"
)

type DispatchInfo struct {
	MetadataPath Path
	DataPath     Path
	ProcessedBy  []ProcessorId
}

type Dispatcher struct {
	workPath      Path
	dispatchChan  chan DispatchInfo
	ingestDirChan chan IngestDirInfo
	processorReg  *ProcessorRegistry
	logger        *log.Logger
}

func NewDispatcher(workPath Path, dispatchChan chan DispatchInfo, ingestDirChan chan IngestDirInfo, processorReg *ProcessorRegistry) (*Dispatcher, error) {
	if !workPath.Exists() {
		return nil, fmt.Errorf("Work path does not exist: %s", workPath)
	}

	logger := log.New(os.Stdout, "Dispatcher: ", 0)
	return &Dispatcher{workPath, dispatchChan, ingestDirChan, processorReg, logger}, nil
}

func (dispatcher *Dispatcher) Run() {
	for dispatchInfo := range dispatcher.dispatchChan {
		metadata, err := PropertiesFromFile(dispatchInfo.MetadataPath)

		if err != nil {
			dispatcher.logger.Print(err.Error())
			continue
		}

		matchingProcessors := dispatcher.processorReg.MatchingProcessors(metadata)

		for _, processor := range matchingProcessors {
			inputPath, err := makeRandomPath(dispatcher.workPath).Resolve()

			if err != nil {
				dispatcher.logger.Print("Error creating input path for dispatch: %s", err.Error())
				continue
			}

			outputPath, err := makeRandomPath(dispatcher.workPath).Resolve()

			if err != nil {
				dispatcher.logger.Print("Error creating output path for dispatch: %s", err.Error())
				continue
			}

			// copy metadata into the new input path
			err = CopyFile(dispatchInfo.MetadataPath, inputPath)

			if err != nil {
				dispatcher.logger.Print("Error copying metadata to work path: %s", err.Error())
				continue
			}

			// copy data into the new input path
			err = CopyFile(dispatchInfo.DataPath, inputPath)
			if err != nil {
				dispatcher.logger.Print("Error copying data to work path: %s", err.Error())
				continue
			}

			// copy the ProcessedBy slice so multiple processors
			// don't update the same slice
			processedByCopy := make([]ProcessorId, len(dispatchInfo.ProcessedBy))
			copy(processedByCopy, dispatchInfo.ProcessedBy)

			go processor.Run(inputPath, outputPath, processedByCopy, dispatcher.ingestDirChan)
		}
	}
}
