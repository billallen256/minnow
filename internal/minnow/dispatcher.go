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

func (info DispatchInfo) AlreadyProcessedBy(processorId ProcessorId) bool {
	for _, id := range info.ProcessedBy {
		if id == processorId {
			return true
		}
	}

	return false
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
			if dispatchInfo.AlreadyProcessedBy(processor.GetId()) {
				dispatcher.logger.Printf("Data at %s already processed by processor %s. Will not process again.", dispatchInfo.DataPath, processor.GetId())
				continue
			}

			inputPath, err := makeRandomPath(dispatcher.workPath)

			if err != nil {
				dispatcher.logger.Printf("Error creating input path for dispatch: %s", err.Error())
				continue
			}

			outputPath, err := makeRandomPath(dispatcher.workPath)

			if err != nil {
				dispatcher.logger.Printf("Error creating output path for dispatch: %s", err.Error())
				continue
			}

			// Need to resolve the input and output paths so
			// clean, absolute paths get passed to the processor.
			inputPath, err = inputPath.Resolve()

			if err != nil {
				dispatcher.logger.Print(err.Error())
				continue
			}

			outputPath, err = outputPath.Resolve()

			if err != nil {
				dispatcher.logger.Print(err.Error())
				continue
			}

			// copy metadata into the new input path
			err = CopyFile(dispatchInfo.MetadataPath, inputPath)

			if err != nil {
				dispatcher.logger.Printf("Error copying metadata to work path: %s", err.Error())
				continue
			}

			// copy data into the new input path
			err = CopyFile(dispatchInfo.DataPath, inputPath)
			if err != nil {
				dispatcher.logger.Printf("Error copying data to work path: %s", err.Error())
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
