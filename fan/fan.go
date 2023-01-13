package fan

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type msg[T any] struct {
	count uint64
	data  T
}

type sub[T any] struct {
	ch chan T
	ls *list.List

	ctx context.Context
	cn  context.CancelFunc

	mu sync.RWMutex
}

func NewSub[T any](ch chan T) *sub[T] {
	o := &sub[T]{}
	o.ch = ch
	o.ls = list.New()
	o.ctx, o.cn = context.WithCancel(context.Background())
	go o.loop()
	return o
}

func (s *sub[T]) loop() {
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-s.ctx.Done():
			ticker.Stop()
		case <-ticker.C:
			s.mu.Lock()
			for s.tick() {
			}
			s.mu.Unlock()
		}
	}
}
func (s *sub[T]) tick() bool {
	// see if can send
	last := s.ls.Back()
	if last == nil {
		return false
	}
	select {
	case s.ch <- last.Value.(T):
		s.ls.Remove(last)
		return true
	case <-s.ctx.Done():
	default:
	}
	return false
}
func (s *sub[T]) fwd() {
	// see if can send
	last := s.ls.Back()
	if last == nil {
		return
	}
	select {
	case s.ch <- last.Value.(T):
		s.ls.Remove(last)
	case <-s.ctx.Done():
		return
	default:
	}
}
func (s *sub[T]) Send(x T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ls.PushFront(x)
	// try to forward its own message before existing
	// if it returns, either its message was sent, or it was closed.
	s.fwd()
}
func (s *sub[T]) Close() {
	select {
	case <-s.ctx.Done():
		return
	default:
	}
	s.cn()
}
