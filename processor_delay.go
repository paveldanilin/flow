package flow

import (
	"context"
	"time"
)

type DelayProcessor struct {
	dur time.Duration
	BaseProcessor
}

func NewDelayProcessor(dur time.Duration) *DelayProcessor {
	return &DelayProcessor{dur: dur}
}

func (p *DelayProcessor) Process(ctx context.Context, exchange *Exchange) error {
	time.Sleep(p.dur)

	return p.next(ctx, exchange)
}
