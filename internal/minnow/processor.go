package minnow

type ProcessorId Path

type Processor struct {
	definitionPath Path
	configPath     Path
	hook           Hook
}

func (processor Processor) GetId() ProcessorId {
	return ProcessorId(pi.definitionPath)
}
