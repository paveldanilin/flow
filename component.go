package flow

type Component interface {
	GetConsumer(params map[string]any) Consumer
	GetProducer(params map[string]any) Producer
}
