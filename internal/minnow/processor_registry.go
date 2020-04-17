package minnow

import (
	"log"
	"sync"
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

	// Start a goroutine to periodically scan for added/removed processors.
	go func(registry *ProcessorRegistry) {
		for range time.Tick(time.Duration(5) * time.Minute) {
			err := registry.BuildProcessorMap()

			if err != nil {
				registry.logger.Print(err.Error())
			}
		}
	}(registry)

	return registry
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
			definitionDirPaths = append(definitionDirPaths, definitionPaths)
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
}

func (registry *ProcessorRegistry) MatchingProcessors(metadata Properties) []*Processor {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()

}
