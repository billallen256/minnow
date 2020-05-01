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

type RunRequest struct {
	inputPath Path
	outputPath Path
	processedBy []ProcessorId
	ingestDirChan chan IngestDirInfo
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
	processorRegistry  *ProcessorRegistry
	logger        *log.Logger
}

func NewDispatcher(workPath Path, dispatchChan chan DispatchInfo, ingestDirChan chan IngestDirInfo, processorRegistry *ProcessorRegistry) (*Dispatcher, error) {
	if !workPath.Exists() {
		return nil, fmt.Errorf("Work path does not exist: %s", workPath)
	}

	logger := log.New(os.Stdout, "Dispatcher: ", 0)
	return &Dispatcher{workPath, dispatchChan, ingestDirChan, processorRegistry, logger}, nil
}

func (dispatcher *Dispatcher) Run() {
	for dispatchInfo := range dispatcher.dispatchChan {
		metadata, err := PropertiesFromFile(dispatchInfo.MetadataPath)

		if err != nil {
			dispatcher.logger.Print(err.Error())
			continue
		}

		matchingProcessorIds := dispatcher.processorRegistry.MatchingProcessorIds(metadata)

		for _, processorId := range matchingProcessorIds {
			if dispatchInfo.AlreadyProcessedBy(processorId) {
				dispatcher.logger.Printf("Data at %s already processed by processor %s. Will not process again.", dispatchInfo.DataPath, processorId)
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
			// clean, _absolute_ paths get passed to the processor.
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

			runRequest := RunRequest{inputPath, outputPath, processedByCopy, dispatcher.ingestDirChan}
			dispatcher.processorRegistry.SendToProcessorId(processorId, runRequest)
		}
	}
}
