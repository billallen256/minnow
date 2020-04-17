package minnow

import (
	"fmt"
	"log"
)

type ProcessorId Path

type Processor struct {
	name           string
	definitionPath Path
	configPath     Path
	hook           Hook
	logger         *log.Logger
}

func NewProcessor(definitionPath Path) (Processor, error) {
	if !definitionPath.IsDir() {
		return Processor{}, fmt.Errorf("Processor definition path must be a directory: %s", definitionPath)
	}

	configPath, err := definitionPath.JoinPath("config.properties")

	if err != nil {
		return Processor{}, err
	}

	if !configPath.Exists() {
		return Processor{}, fmt.Errorf("Processor config file does not exist: %s", configPath)
	}

}

func (processor Processor) GetId() ProcessorId {
	return ProcessorId(processor.definitionPath)
}

func (processor Processor) Run(inputPath, outputPath Path, ingestChan chan<- IngestInfo) error {

}
