package minnow

import (
	"log"
	"sync"
)

type ProcessorRegistry struct {
	DefinitionsPath Path
	Processors      map[ProcessorId]Processor
	Mutex           *sync.RWMutex
	Logger          *log.Logger
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
				registry.Logger.Print(err.Error())
			}
		}
	}(registry)

	return registry
}

func (registry *ProcessorRegistry) BuildProcessorMap() error {

}

func (registry *ProcessorRegistry) MatchingProcessors(metadata Properties) []*Processor {

}
