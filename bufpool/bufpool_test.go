package bufpool

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
)

func TestAllocations(t *testing.T) {
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	runtime.GC()
	for i := 0; i < 10000; i++ {
		b := Get(1010)
		Put(b)
	}
	runtime.GC()
	runtime.ReadMemStats(&m2)
	frees := m2.Frees - m1.Frees
	if frees > 1000 {
		t.Fatalf("expected less than 100 frees after GC, got %d", frees)
	}
}

func TestRange(t *testing.T) {
	min := nextLogBase2(1)
	max := nextLogBase2(uint32(MaxLength))
	if int(max) != len(GlobalPool.pools)-1 {
		t.Errorf("expected %d pools, found %d", max, len(GlobalPool.pools))
	}
	if min != 0 {
		t.Errorf("unused min pool")
	}
}

func TestPool(t *testing.T) {
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var p BufferPool

	a := make([]byte, 21)
	a[0] = 1
	b := make([]byte, 2050)
	b[0] = 2
	ab := bytes.NewBuffer(a)
	bb := bytes.NewBuffer(b)
	p.Put(ab)
	p.Put(bb)
	if g := p.Get(16); g.Cap() > 21 {
		t.Fatalf("got [%d,...]; want [1,...]", g.Cap())
	}
	if g := p.Get(2048); g.Cap() < 2048 {
		t.Fatalf("got [%d,...]; want [2,...]", g.Cap())
	}
	if g := p.Get(16); g.Cap() != 16 {
		t.Fatalf("got existing slice; want new slice")
	}
	if g := p.Get(2048); g.Cap() != 2048 {
		t.Fatalf("got existing slice; want new slice")
	}
	if g := p.Get(1); g.Cap() != 1 {
		t.Fatalf("got existing slice; want new slice")
	}
	d := make([]byte, 1023)
	d[0] = 3
	p.Put(bytes.NewBuffer(d))
	if g := p.Get(1024); g.Cap() != 1024 {
		t.Fatalf("got existing slice; want new slice")
	}
	if g := p.Get(512); g.Cap() != 1023 {
		t.Fatalf("got [%v,...]; want [3,...]", "TODO: fix test")
	}
	p.Put(ab)

	debug.SetGCPercent(100) // to allow following GC to actually run
	runtime.GC()
	// For some reason, you need to run GC twice on go 1.16 if you want it to reliably work.
	runtime.GC()
	p.Get(10)
}

func TestPoolStressByteSlicePool(t *testing.T) {
	var p BufferPool
	const P = 10
	chs := 10
	maxSize := 1 << 16
	N := int(1e4)
	if testing.Short() {
		N /= 100
	}
	done := make(chan bool)
	errs := make(chan error)
	for i := 0; i < P; i++ {
		go func() {
			ch := make(chan *bytes.Buffer, chs+1)

			for i := 0; i < chs; i++ {
				j := rand.Int() % maxSize
				ch <- p.Get(j)
			}

			for j := 0; j < N; j++ {
				r := 0
				for i := 0; i < chs; i++ {
					v := <-ch
					p.Put(v)
					r = rand.Int() % maxSize
					v = p.Get(r)
					if v.Cap() < r {
						errs <- fmt.Errorf("expect len(v) >= %d, got %d", j, v.Len())
					}
					ch <- v
				}

				if r%1000 == 0 {
					runtime.GC()
				}
			}
			done <- true
		}()
	}

	for i := 0; i < P; {
		select {
		case <-done:
			i++
		case err := <-errs:
			t.Error(err)
		}
	}
}
func BenchmarkStdPoolFlatDistribution(b *testing.B) {
	var p = sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 7
		for pb.Next() {
			if i > 1<<20 {
				i = 7
			} else {
				i = i << 1
			}
			b := p.Get().(*bytes.Buffer)
			b.WriteByte(byte(i))
			b.Reset()
			p.Put(b)
		}
	})
}

func BenchmarkPool(b *testing.B) {
	var p BufferPool
	b.RunParallel(func(pb *testing.PB) {
		i := 7
		for pb.Next() {
			if i > 1<<20 {
				i = 7
			} else {
				i = i << 1
			}
			b := p.Get(i)
			b.WriteByte(byte(i))
			p.Put(b)
		}
	})
}
func BenchmarkPoolStd(b *testing.B) {
	var p BufferPool
	b.RunParallel(func(pb *testing.PB) {
		i := 7
		for pb.Next() {
			if i > 1<<20 {
				i = 7
			} else {
				i = i << 1
			}
			b := p.GetStd()
			b.WriteByte(byte(i))
			p.PutStd(b)
		}
	})
}

func BenchmarkAlloc(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		i := 7
		for pb.Next() {
			if i > 1<<20 {
				i = 7
			} else {
				i = i << 1
			}
			b := make([]byte, i)
			b[1] = byte(i)
		}
	})
}

func BenchmarkPoolOverlflow(b *testing.B) {
	var p BufferPool
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bufs := make([][]byte, 2100)
			for pow := uint32(0); pow < 21; pow++ {
				for i := 0; i < 100; i++ {
					bufs = append(bufs, p.Get(1<<pow).Bytes())
				}
			}
			for _, b := range bufs {
				p.Put(bytes.NewBuffer(b))
			}
		}
	})
}

func ExampleGet() {
	buf := Get(100)
	fmt.Println("capacity", buf.Cap())
	Put(buf)
	// Output:
	// capacity 128
}
