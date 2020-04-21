package minnow

import (
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	IngestPath               Path
	IngestMinAge             time.Duration
	WorkPath                 Path
	WorkAgeOff               time.Duration
	ProcessorDefinitionsPath Path
}

func ReadConfig(path Path) (Config, error) {
	configProperties, err := PropertiesFromFile(path)

	if err != nil {
		return Config{}, err
	}

	return ParseConfig(configProperties)
}

func ParseConfig(configProperties Properties) (Config, error) {
	config := Config{}

	//
	// ingest_path
	//
	ingestPathStr, found := configProperties["ingest_dir"]

	if !found {
		return Config{}, fmt.Errorf("ingest_dir missing from config file")
	}

	config.IngestPath = Path(ingestPathStr)

	if !config.IngestPath.Exists() {
		return Config{}, fmt.Errorf("ingest_dir does not exist at %s", config.IngestPath)
	}

	//
	// ingest_min_age
	//
	ingestMinAgeStr, found := configProperties["ingest_min_age"]

	if !found {
		ingestMinAgeStr = "300" // set default of five minutes
	}

	ingestMinAgeInt, err := strconv.Atoi(ingestMinAgeStr)

	if err != nil {
		return Config{}, fmt.Errorf("ingest_min_age must be an integer representing seconds")
	}

	config.IngestMinAge = time.Duration(ingestMinAgeInt) * time.Second

	//
	// work_path
	//
	workPathStr, found := configProperties["work_dir"]

	if !found {
		return Config{}, fmt.Errorf("work_dir missing from config file")
	}

	config.WorkPath = Path(workPathStr)

	if !config.WorkPath.Exists() {
		return Config{}, fmt.Errorf("work_dir does not exist at %s", config.WorkPath)
	}

	//
	// work_age_off
	//
	workAgeOffStr, found := configProperties["work_age_off"]

	if !found {
		workAgeOffStr = "172800" // set default of two days
	}

	workAgeOffInt, err := strconv.Atoi(workAgeOffStr)

	if err != nil {
		return Config{}, fmt.Errorf("work_age_off must be an integer representing seconds")
	}

	config.WorkAgeOff = time.Duration(workAgeOffInt) * time.Second

	//
	// processor_definitions_dir
	//
	processorDefinitionsPathStr, found := configProperties["processor_definitions_dir"]

	if !found {
		return Config{}, fmt.Errorf("processor_definitions_dir missing from config file")
	}

	config.ProcessorDefinitionsPath = Path(processorDefinitionsPathStr)

	if !config.ProcessorDefinitionsPath.Exists() {
		return Config{}, fmt.Errorf("processor_definitions_path does not exist at %s", config.ProcessorDefinitionsPath)
	}

	return config, nil
}
