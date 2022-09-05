package bufpool

import (
	"bytes"
	"math"
	"math/bits"
	"sync"
)

const MaxLength = math.MaxInt32

type BufferPool struct {
	pools [32]sync.Pool
}

var GlobalPool = new(BufferPool)

var allocator = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// Get retrieves a buffer of the appropriate length from the buffer pool or
// allocates a new one. Get may choose to ignore the pool and treat it as empty.
// Callers should not assume any relation between values passed to Put and the
// values returned by Get.
//
// If no suitable buffer exists in the pool, Get creates one.
func (p *BufferPool) GetStd() *bytes.Buffer {
	return allocator.Get().(*bytes.Buffer)
}

func (p *BufferPool) PutStd(b *bytes.Buffer) {
	allocator.Put(b)
}
func (p *BufferPool) Get(length int) *bytes.Buffer {
	if length == 0 {
		return allocator.Get().(*bytes.Buffer)
	}
	// Calling this function with a negative length is invalid.
	// make will panic if length is negative, so we don't have to.
	if length > MaxLength || length < 0 {
		return allocator.Get().(*bytes.Buffer)
	}
	idx := nextLogBase2(uint32(length))
	if ptr := p.pools[idx].Get(); ptr != nil {
		bp := ptr.(*bytes.Buffer)
		bp.Reset()
		return bp
	}
	return bytes.NewBuffer(make([]byte, 0, 1<<idx)[:uint32(length)])
}

func (p *BufferPool) New(length int) *bytes.Buffer {
	return p.Get(length)
}

// Put adds x to the pool.
func (p *BufferPool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	capacity := buf.Cap()
	if capacity == 0 || capacity > MaxLength {
		allocator.Put(buf)
		return // drop it
	}
	idx := prevLogBase2(uint32(capacity))
	p.pools[idx].Put(buf)
}

// Get retrieves a buffer of the appropriate length from the global buffer pool
// (or allocates a new one).
func Get(length int) *bytes.Buffer {
	return GlobalPool.Get(length)
}

// Put returns a buffer to the global buffer pool.
func Put(slice *bytes.Buffer) {
	GlobalPool.Put(slice)
}

// Log of base two, round up (for v > 0).
func nextLogBase2(v uint32) uint32 {
	return uint32(bits.Len32(v - 1))
}

// Log of base two, round down (for v > 0)
func prevLogBase2(num uint32) uint32 {
	next := nextLogBase2(num)
	if num == (1 << uint32(next)) {
		return next
	}
	return next - 1
}
