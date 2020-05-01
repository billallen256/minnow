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
	processorPools      map[ProcessorId]*ProcessorPool
	mutex           *sync.RWMutex
	logger          *log.Logger
}

func NewProcessorRegistry(definitionsPath Path) (*ProcessorRegistry, error) {
	processorPools := make(map[ProcessorId]*ProcessorPool)
	mutex := new(sync.RWMutex)
	logger := log.New(os.Stdout, "ProcessorRegistry: ", 0)
	registry := &ProcessorRegistry{definitionsPath, processorPools, mutex, logger}

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

	processorPools := make(map[ProcessorId]*ProcessorPool)

	for _, definitionDirPath := range definitionDirPaths {
		processor, err := NewProcessor(definitionDirPath)

		if err != nil {
			registry.logger.Print(err.Error())
			continue
		}

		processorPool := NewProcessorPool(processor, processor.GetPoolSize())
		processorPools[processor.GetId()] = processorPool
	}

	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	// stop old ProcessorPool queues so they don't get more data
	for _, processorPool := range registry.processorPools {
		processorPool.Stop()
	}

	registry.processorPools = processorPools

	if len(processorPools) == 0 {
		return fmt.Errorf("No processors found in %s", registry.definitionsPath)
	}

	return nil
}

func (registry *ProcessorRegistry) MatchingProcessorIds(metadata Properties) []ProcessorId {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()

	matchingProcessorIds := make([]ProcessorId, 0)

	for _, processorPool := range registry.processorPools {
		if processorPool.ProcessorHookMatches(metadata) {
			matchingProcessorIds = append(matchingProcessorIds, processorPool.GetProcessorId())
		}
	}

	return matchingProcessorIds
}

func (registry *ProcessorRegistry) SendToProcessorId(processorId ProcessorId, runRequest RunRequest) error {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()

	if processorPool, found := registry.processorPools[processorId]; found {
		processorPool.Run(runRequest)
		return nil
	}

	return fmt.Errorf("Could not send RunRequest to ProcessorId %s", processorId)
}
