package minnow

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type ProcessorId Path

type Processor struct {
	name           string
	definitionPath Path
	executable     string
	hook           Hook
	poolSize       int
	logger         *log.Logger
}

type ProcessorConfig struct {
	Executable string
	Hook       Hook
	PoolSize   int
}

func NewProcessor(definitionPath Path) (Processor, error) {
	if !definitionPath.IsDir() {
		return Processor{}, fmt.Errorf("Processor definition path must be a directory: %s", definitionPath)
	}

	configPath := definitionPath.JoinPath("config.properties")

	if !configPath.Exists() {
		return Processor{}, fmt.Errorf("Processor config file does not exist: %s", configPath)
	}

	config, err := parseProcessorConfig(configPath, definitionPath)

	if err != nil {
		return Processor{}, err
	}

	name := definitionPath.Name()
	logger := log.New(os.Stdout, name+": ", 0)
	return Processor{name, definitionPath, config.Executable, config.Hook, config.PoolSize, logger}, nil
}

func parseProcessorConfig(configPath, definitionPath Path) (ProcessorConfig, error) {
	configProperties, err := PropertiesFromFile(configPath)

	if err != nil {
		return ProcessorConfig{}, err
	}

	executable, found := configProperties["executable"]

	if !found {
		return ProcessorConfig{}, fmt.Errorf("Processor config missing executable property")
	}

	executablePath := definitionPath.JoinPath(Path(executable))

	if !executablePath.Exists() {
		return ProcessorConfig{}, fmt.Errorf("Could not find executable at %s", executablePath)
	}

	poolSizeString, found := configProperties["pool_size"]

	if !found {
		poolSizeString = "5"
	}

	poolSize, err := strconv.Atoi(poolSizeString)

	if err != nil {
		return ProcessorConfig{}, fmt.Errorf("Invalid value for pool_size: %s", err.Error())
	}

	if poolSize < 0 {
		return ProcessorConfig{}, fmt.Errorf("pool_size must be a positive value")
	}

	hookPathString, found := configProperties["hook_file"]

	if !found {
		return ProcessorConfig{}, fmt.Errorf("Processor config missing hook_file property")
	}

	hookPath, err := definitionPath.JoinPath(Path(hookPathString)).Resolve()

	if err != nil {
		return ProcessorConfig{}, err
	}

	// handle hook_type last since it defaults to fail
	hookType, found := configProperties["hook_type"]

	if !found {
		hookType = "basicpropertiesmatchhook"
	}

	if hookType == "basicpropertiesmatchhook" {
		hook, err := NewBasicPropertiesMatchHookFromFile(hookPath)

		if err != nil {
			return ProcessorConfig{}, err
		}

		return ProcessorConfig{executable, hook, poolSize}, nil
	}

	return ProcessorConfig{}, fmt.Errorf("Unknown hook_type %s", hookType)
}

func (processor Processor) GetId() ProcessorId {
	return ProcessorId(processor.definitionPath)
}

func (processor Processor) GetName() string {
	return processor.name
}

func (processor Processor) GetPoolSize() int {
	return processor.poolSize
}

func (processor Processor) Run(runRequestQueue chan RunRequest) {
	for runRequest := range runRequestQueue {
		processor.RunCommand(runRequest)
	}
}

func (processor Processor) RunCommand(runRequest RunRequest) error {
	cmd := exec.Command("./"+processor.executable, string(runRequest.InputPath), string(runRequest.OutputPath))
	cmd.Dir = string(processor.definitionPath) // set the working directory for the command
	processor.logger.Printf("Processor %s running %s", processor.name, cmd.String())
	stdoutStderr, err := cmd.CombinedOutput()
	processorOutputPath := runRequest.OutputPath.JoinPath(Path(fmt.Sprintf("_%s_output.txt", processor.name)))
	outputErr := processorOutputPath.WriteBytes(stdoutStderr)

	if outputErr != nil {
		processor.logger.Printf("Processor %s could not write stdout/stderr to %s: %s", processor.name, processorOutputPath, outputErr.Error())
		// don't return the error here, just log it
	}

	if err != nil {
		processor.logger.Printf("Processor %s returned error: %s", processor.name, err.Error())
		return err
	}

	processor.logger.Print("Processor completed successfully")
	processedBy := append(runRequest.ProcessedBy, processor.GetId())
	runRequest.IngestDirChan <- IngestDirInfo{runRequest.OutputPath, time.Duration(0), processedBy, true}
	err = runRequest.InputPath.RmdirRecursive() // make sure the input directory gets removed

	if err != nil {
		processor.logger.Printf("Could not remove input path: %s", err.Error())
	}

	return nil
}

func (processor Processor) HookMatches(properties Properties) bool {
	return processor.hook.Matches(properties)
}
