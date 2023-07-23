package flow

import "time"

type DelayProcessor struct {
	dur time.Duration
	BaseProcessor
}

func NewDelayProcessor(dur time.Duration) *DelayProcessor {
	return &DelayProcessor{dur: dur}
}

func (p *DelayProcessor) Process(exchange *Exchange) error {
	time.Sleep(p.dur)

	return p.next(exchange)
}
