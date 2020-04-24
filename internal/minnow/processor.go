package minnow

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type ProcessorId Path

type Processor struct {
	name           string
	definitionPath Path
	executable     string
	hook           Hook
	logger         *log.Logger
}

type ProcessorConfig struct {
	Executable string
	Hook       Hook
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
	return Processor{name, definitionPath, config.Executable, config.Hook, logger}, nil
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

	hookPathString, found := configProperties["hook_file"]

	if !found {
		return ProcessorConfig{}, fmt.Errorf("Processor config missing hook_file property")
	}

	hookPath, err := definitionPath.JoinPath(Path(hookPathString)).Resolve()

	if err != nil {
		return ProcessorConfig{}, err
	}

	hookType, found := configProperties["hook_type"]

	if !found {
		hookType = "basicpropertiesmatchhook"
	}

	if hookType == "basicpropertiesmatchhook" {
		hook, err := NewBasicPropertiesMatchHookFromFile(hookPath)

		if err != nil {
			return ProcessorConfig{}, err
		}

		return ProcessorConfig{executable, hook}, nil
	}

	return ProcessorConfig{}, fmt.Errorf("Unknown hook_type %s", hookType)
}

func (processor Processor) GetId() ProcessorId {
	return ProcessorId(processor.definitionPath)
}

func (processor Processor) Run(inputPath, outputPath Path, processedBy []ProcessorId, ingestDirChan chan<- IngestDirInfo) error {
	cmd := exec.Command("./"+processor.executable, string(inputPath), string(outputPath))
	cmd.Dir = string(processor.definitionPath) // set the working directory for the command
	processor.logger.Printf("Processor %s running %s", processor.name, cmd.String())
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		processor.logger.Printf("Processor %s returned error: %s", processor.name, err.Error())
		return err
	}

	processor.logger.Print("Processor completed successfully")
	processorOutputPath := outputPath.JoinPath("processor_output.txt")
	err = processorOutputPath.WriteBytes(stdoutStderr)

	if err != nil {
		processor.logger.Printf("Processor %s could not write output to %s", processor.name, processorOutputPath)
		// don't return the error here, just log it
	}

	ingestDirChan <- IngestDirInfo{outputPath, time.Duration(0), append(processedBy, processor.GetId())}
	return nil
}

func (processor Processor) HookMatches(properties Properties) bool {
	return processor.hook.Matches(properties)
}
