package minnow

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type ProcessorRegistry struct {
	definitionsPath Path
	processors      map[ProcessorId]Processor
	mutex           *sync.RWMutex
	logger          *log.Logger
}

func NewProcessorRegistry(definitionsPath Path) (*ProcessorRegistry, error) {
	processors := make(map[ProcessorId]Processor)
	mutex := new(sync.RWMutex)
	logger := log.New(os.Stdout, "ProcessorRegistry: ", 0)
	registry := &ProcessorRegistry{definitionsPath, processors, mutex, logger}

	err := registry.BuildProcessorMap()

	if err != nil {
		return nil, err
	}

	return registry, nil
}

func (registry *ProcessorRegistry) Run() {
	for range time.Tick(time.Duration(5) * time.Minute) {
		err := registry.BuildProcessorMap()

		if err != nil {
			registry.logger.Print(err.Error())
		}
	}
}

func (registry *ProcessorRegistry) BuildProcessorMap() error {
	definitionPaths, err := registry.definitionsPath.Glob("*")

	if err != nil {
		return err
	}

	// Filter out any non-directories
	definitionDirPaths := make([]Path, 0)

	for _, definitionPath := range definitionPaths {
		if definitionPath.IsDir() {
			definitionDirPaths = append(definitionDirPaths, definitionPath)
		}
	}

	processors := make(map[ProcessorId]Processor)

	for _, definitionDirPath := range definitionDirPaths {
		processor, err := NewProcessor(definitionDirPath)

		if err != nil {
			registry.logger.Print(err.Error())
			continue
		}

		processors[processor.GetId()] = processor
	}

	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	registry.processors = processors

	if len(processors) == 0 {
		return fmt.Errorf("No processors found in %s", registry.definitionsPath)
	}

	return nil
}

func (registry *ProcessorRegistry) MatchingProcessors(metadata Properties) []Processor {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()

	matchingProcessors := make([]Processor, 0)

	for _, processor := range registry.processors {
		if processor.HookMatches(metadata) {
			matchingProcessors = append(matchingProcessors, processor)
		}
	}

	return matchingProcessors
}
