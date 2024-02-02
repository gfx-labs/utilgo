package mint_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"gfx.cafe/util/go/mint"
	"gfx.cafe/util/go/mint/hub"
)

type event struct {
	F1, F2 string
}

func TestEmitSimple(t *testing.T) {
	e := new(mint.Emitter)

	received := false
	off := mint.On(e, func(ctx context.Context, e event) { received = true })
	defer off()

	mint.Emit(nil, e, event{"hello", "world"})

	if !received {
		t.Fatalf("didn't receive")
	}
}

func TestEmitRecursive(t *testing.T) {
	e := new(mint.Emitter)

	var i int
	mint.On(e, func(context.Context, event) {
		if i < 5 {
			i += 1
			mint.Emit(nil, e, event{})
		}
	})

	mint.Emit(nil, e, event{})

	if i != 5 {
		t.Fatalf("didn't receive")
	}
}

func TestEmitConcurrent(t *testing.T) {
	e := new(mint.Emitter)

	var i atomic.Uint32
	mint.On(e, func(context.Context, event) {
		i.Add(1)
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			mint.Emit(nil, e, event{})
			wg.Done()
		}()
	}

	wg.Wait()
	if i := i.Load(); i != 100 {
		t.Fatalf("lost emits; got %d expected 100", i)
	}
}

func TestBroadReceiverMisfire(t *testing.T) {
	e := new(mint.Emitter)

	mint.On(e, func(context.Context, any) { t.Error("misfired 'any' consumer with 'event' emit") })
	mint.Emit(nil, e, event{})
}

func TestOffSimple(t *testing.T) {
	e := new(mint.Emitter)

	c := 0
	off := mint.On(e, func(_ context.Context, v int) { c = v })

	mint.Emit(nil, e, 1)
	<-off() // wait for it to synchronize
	mint.Emit(nil, e, 2)

	if c != 1 {
		t.Fatalf("expected c to be %d; got %d", 1, c)
	}
}

func TestContextCancel(t *testing.T) {
	e := new(mint.Emitter)

	hub.On(e, func(_ context.Context, v event) {
		t.Errorf("consumer called despite context cancel")
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := hub.Emit(e, ctx, event{})
	if err != context.Canceled {
		t.Errorf("expected context.Canceled; got %v", err)
	}
}

func TestContextNoEmitter(t *testing.T) {
	ctx := context.Background()
	if err := hub.Emit(nil, ctx, event{}); err != nil {
		t.Errorf("expected error; got %v", err)
	}
}

func TestUse(t *testing.T) {
	e := new(mint.Emitter)

	s := make([]int, 0, 3)
	mint.Use(e, func(context.Context, any) func() {
		s = append(s, 1)
		return func() { s = append(s, 3) }
	})

	mint.On(e, func(_ context.Context, e event) { s = append(s, 2) })

	mint.Emit(nil, e, event{})

	if len(s) != 3 || s[0] != 1 || s[1] != 2 || s[2] != 3 {
		t.Fatalf("plugin was called in incorrect order: expected %v, got %v", []int{1, 2, 3}, s)
	}
}
