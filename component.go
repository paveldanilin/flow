package flow

import (
	"context"
	"errors"
	"github.com/paveldanilin/flow/uri"
	"sync"
)

type Daemon interface {
	Start() error
	Stop()
}

type Consumer interface {
	Start() error
	Stop()
}

type Producer interface {
	Process(ctx context.Context, exchange *Exchange) error
}

type Component interface {
	GetConsumer(componentURI string) (Consumer, error)
	GetProducer(componentURI string) (Producer, error)
}

var componentMap = map[string]Component{}
var componentMu = sync.RWMutex{}

func RegisterComponent(name string, component Component) error {
	componentMu.Lock()
	defer componentMu.Unlock()

	if _, componentExists := componentMap[name]; componentExists {
		return errors.New("flow: component already exists")
	}

	componentMap[name] = component

	return nil
}

func startComponents() error {
	componentMu.RLock()
	defer componentMu.RUnlock()

	for _, comp := range componentMap {
		if daemon, isDeamon := comp.(Daemon); isDeamon {
			err := daemon.Start()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func stopComponents() {
	componentMu.RLock()
	defer componentMu.RUnlock()

	for _, comp := range componentMap {
		if daemon, isDeamon := comp.(Daemon); isDeamon {
			daemon.Stop()
		}
	}
}

func getComponent(componentURI string) (Component, error) {
	componentName := uri.Schema(componentURI)
	componentMu.RLock()
	component, exists := componentMap[componentName]
	componentMu.RUnlock()

	if !exists {
		return nil, errors.New("unknown component")
	}

	return component, nil
}

func getConsumer(consumerURI string) (Consumer, error) {
	component, err := getComponent(consumerURI)
	if err != nil {
		return nil, err
	}

	return component.GetConsumer(consumerURI)
}

func getProducer(producerURI string) (Producer, error) {
	component, err := getComponent(producerURI)
	if err != nil {
		return nil, err
	}

	return component.GetProducer(producerURI)
}
