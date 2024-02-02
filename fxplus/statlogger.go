package fxplus

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/c2h5oh/datasize"
)

func StatLogger(ctx context.Context, l *slog.Logger) {
	rtm := runtime.MemStats{
		Alloc:         0,
		TotalAlloc:    0,
		Sys:           0,
		Lookups:       0,
		Mallocs:       0,
		Frees:         0,
		HeapAlloc:     0,
		HeapSys:       0,
		HeapIdle:      0,
		HeapInuse:     0,
		HeapReleased:  0,
		HeapObjects:   0,
		StackInuse:    0,
		StackSys:      0,
		MSpanInuse:    0,
		MSpanSys:      0,
		MCacheInuse:   0,
		MCacheSys:     0,
		BuckHashSys:   0,
		GCSys:         0,
		OtherSys:      0,
		NextGC:        0,
		LastGC:        0,
		PauseTotalNs:  0,
		NumGC:         0,
		NumForcedGC:   0,
		GCCPUFraction: 0.0,
	}
	interval := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-interval.C:
			case <-ctx.Done():
				return
			}
			runtime.ReadMemStats(&rtm)
			l.Debug("runtime",
				"sys", datasize.ByteSize(rtm.Sys).HR(),
				"halloc", datasize.ByteSize(rtm.HeapAlloc).HR(),
				"hsys", datasize.ByteSize(rtm.HeapSys).HR(),
				"hidle", datasize.ByteSize(rtm.HeapIdle).HR(),
				"hinuse", datasize.ByteSize(rtm.HeapInuse).HR(),
				"hrelease", datasize.ByteSize(rtm.HeapReleased).HR(),
				"sinuse", datasize.ByteSize(rtm.StackInuse).HR(),
				"ssys", datasize.ByteSize(rtm.StackSys).HR(),
				"NumGc", int(rtm.NumGC),
				"NextGc", datasize.ByteSize(rtm.NextGC).HR(),
				"pause ns", rtm.PauseTotalNs,
			)
		}
	}()
}
