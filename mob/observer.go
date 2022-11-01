package mob

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"gfx.cafe/util/go/generic"
)

// EventHandler provides an interface for an event handler.
type EventHandler[T any] interface {
	Handle(ctx context.Context, event T) error
}

// EventHandlerFunc type is an adapter to allow the use of ordinary functions as event handlers.
type EventHandlerFunc[T any] func(ctx context.Context, event T) error

func (f EventHandlerFunc[T]) Handle(ctx context.Context, event T) error {
	return f(ctx, event)
}

// EventHub is the interface that wraps the mob's Notify method.
type EventHub[T any] interface {
	// Notify dispatches a given event and execute all handlers registered with a dispatched event's type.
	// Handlers are executed concurrently and errors are collected, if any, they're returned to the client.
	//
	// If there is no appropriate handler in the notifier's Mob instance, ErrHandlerNotFound is returned.
	Notify(ctx context.Context, event T) error

	Register(ehn EventHandler[T], opts ...Option)
	RegisterFunc(ehn EventHandlerFunc[T], opts ...Option)
}

// NewEventHub returns an event notifier which uses a given Mob instance.
func NewEventHub[T any](m *Mob) EventHub[T] {
	return &notifier[T]{m: m}
}

func (nf *notifier[T]) Register(evt EventHandler[T], opts ...Option) {
	RegisterEventHandlerTo(nf.m, evt, opts...)
}

func (nf *notifier[T]) RegisterFunc(evt EventHandlerFunc[T], opts ...Option) {
	nf.Register(evt, opts...)
}

// A notifier is a facilitator for a given event type.
type notifier[T any] struct {
	m *Mob

	mu sync.RWMutex
}

func (nf *notifier[T]) Notify(ctx context.Context, event T) error {
	nf.mu.RLock()
	hns, ok := nf.m.ehandlers[reflect.TypeOf(event)]
	if !ok {
		return ErrHandlerNotFound
	}
	nf.mu.RUnlock()
	c := make(chan error)
	var wg sync.WaitGroup
	hns.Range(func(key int64, hn *handler) bool {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Dispatching result not checked because if a handler is found then it should always satisfy EventHandler[T] interface.
			dhn, _ := hn.embedded.(EventHandler[T])
			if err := dhn.Handle(ctx, event); err != nil {
				if hn.name != "" {
					err = fmt.Errorf("%s: %w", hn.name, err)
				}
				hns.Delete(key)
				c <- err
			}
		}()
		return true
	})
	go func() {
		wg.Wait()
		close(c)
	}()
	var aggr AggregateHandlerError = make([]error, 0)
	for err := range c {
		aggr = append(aggr, err)
	}
	if len(aggr) > 0 {
		return aggr
	}
	return nil
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// RegisterEventHandlerTo adds a given event handler to the given Mob instance.
// Returns nil if the handler added successfully, an error otherwise.
//
// Multiple event handlers can be registered for a single event's type.
func RegisterEventHandlerTo[T any](m *Mob, ehn EventHandler[T], opts ...Option) error {
	if !isValid(ehn) {
		return ErrInvalidHandler
	}
	var ev T
	k := reflect.TypeOf(ev)
	hn := &handler{embedded: ehn}
	for _, opt := range opts {
		opt.apply(hn)
	}
	idx := m.eidx.Add(1)
	if _, ok := m.ehandlers[k]; !ok {
		m.ehandlers[k] = &generic.Map[int64, *handler]{}
	}
	m.ehandlers[k].Store(idx, hn)
	return nil
}

// RegisterEventHandlerFuncTo adds a given event handler to the given Mob instance.
// Returns nil if the handler added successfully, an error otherwise.
//
// Multiple event handlers can be registered for a single event's type.
func RegisterEventHandlerFuncTo[T any](m *Mob, ehn EventHandlerFunc[T], opts ...Option) error {
	return RegisterEventHandlerTo[T](m, ehn, opts...)
}

// RegisterEventHandler adds a given event handler to the global Mob instance.
// Returns nil if the handler added successfully, an error otherwise.
//
// Multiple event handlers can be registered for a single event's type.
func RegisterEventHandler[T any](hn EventHandler[T], opts ...Option) error {
	return RegisterEventHandlerTo(m, hn, opts...)
}

// Notify dispatches a given event and execute all handlers registered with a dispatched event's type.
// Handlers are executed concurrently and errors are collected, if any, they're returned to the client.
//
// If there is no appropriate handler in the global Mob instance, ErrHandlerNotFound is returned.
func Notify[T any](ctx context.Context, event T) error {
	return NewEventHub[T](m).Notify(ctx, event)
}
