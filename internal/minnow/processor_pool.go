package minnow

type ProcessorPool struct {
	processor       Processor
	runRequestQueue chan RunRequest
}

func NewProcessorPool(processor Processor, poolSize int) *ProcessorPool {
	runRequestQueue := make(chan RunRequest, 100*poolSize)

	for i := 0; i < poolSize; i++ {
		go processor.Run(runRequestQueue)
	}

	return &ProcessorPool{processor, runRequestQueue}
}

func (pool *ProcessorPool) Run(runRequest RunRequest) {
	pool.runRequestQueue <- runRequest
}

func (pool *ProcessorPool) Stop() {
	close(pool.runRequestQueue)
}

func (pool *ProcessorPool) GetProcessorId() ProcessorId {
	return pool.processor.GetId()
}

func (pool *ProcessorPool) ProcessorHookMatches(properties Properties) bool {
	return pool.processor.HookMatches(properties)
}
