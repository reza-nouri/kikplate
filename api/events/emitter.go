// internal/events/emitter.go
package events

import "sync"

type Handler func(payload any)

type EventEmitter struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		handlers: make(map[string][]Handler),
	}
}

func (e *EventEmitter) On(event string, handler Handler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[event] = append(e.handlers[event], handler)
}

func (e *EventEmitter) Emit(event string, payload any) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, handler := range e.handlers[event] {
		go handler(payload)
	}
}
