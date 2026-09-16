package nkv

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/gfx-labs/utilgo/lambda"
	"github.com/gfx-labs/utilgo/retry"
	"golang.org/x/sync/singleflight"
)

type DB[T any] struct {
	fetcher     Fetcher
	transformer Transformer[T]
	backOff     retry.BackOff

	kvStore map[string]string

	sf singleflight.Group

	ttl         time.Duration
	lastUpdated time.Time
	m           sync.RWMutex
}

// Fetches the nkv formatted file and returns a reader which will be used to parse
// the file. The reader should not be closed by the caller.
type Fetcher func(context.Context) (io.ReadCloser, error)

// Transforms the map[string]string into the desired type. The map is a copy of the
// map returned by the Fetcher.
type Transformer[T any] func(map[string]string) (T, error)

func NewDB[T any](getter Fetcher, transformer Transformer[T], backOff retry.BackOff, ttl time.Duration) *DB[T] {
	return &DB[T]{
		fetcher:     getter,
		transformer: transformer,
		backOff:     backOff,
		ttl:         ttl,
	}
}

func (d *DB[T]) Value(ctx context.Context) (t T, err error) {
	doRefresh := func() (interface{}, error) {
		return d.refresh(ctx)
	}
	if d.ttl >= 0 && time.Since(d.lastUpdated) >= d.ttl || d.lastUpdated == (time.Time{}) {
		v, err, _ := d.sf.Do(strconv.Itoa(int(d.lastUpdated.Unix())), doRefresh)
		if err != nil {
			return t, err
		}
		return v.(T), err
	}
	return d.transform()
}

func (d *DB[T]) transform() (T, error) {
	return d.transformer(d.copy())
}

func (d *DB[T]) copy() map[string]string {
	d.m.Lock()
	defer d.m.Unlock()
	return lambda.MergeMap(make(map[string]string, len(d.kvStore)), d.kvStore)
}

func (d *DB[T]) setState(m map[string]string, updateTime time.Time) map[string]string {
	d.m.Lock()
	defer d.m.Unlock()
	if d.lastUpdated.Before(updateTime) {
		d.kvStore = m
		d.lastUpdated = updateTime
	}
	return lambda.MergeMap(make(map[string]string, len(d.kvStore)), d.kvStore)
}

func (d *DB[T]) refresh(ctx context.Context) (t T, err error) {
	fetchStart := time.Now()
	var reader io.ReadCloser
	err = retry.Until(ctx, d.backOff, func(ctx context.Context, i int) error {
		reader, err = d.fetcher(ctx)
		return err
	})
	if err != nil {
		return
	}
	defer reader.Close()
	responseMap, err := parseIntoMap(reader)
	if err != nil {
		return *new(T), err
	}
	copied := d.setState(responseMap, fetchStart)
	return d.transformer(copied)
}

func parseIntoMap(r io.ReadCloser) (map[string]string, error) {
	scanner := bufio.NewScanner(r)
	parsed := make(map[string]string)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' || bytes.HasPrefix(line, []byte("//")) {
			continue
		}
		k, v, _ := bytes.Cut(line, []byte("="))
		parsed[string(k)] = string(v)
	}
	return parsed, nil
}
