package flow

type Consumer interface {
	Start() error
	Stop()
	Processor
}

type Producer interface {
	Processor
}
