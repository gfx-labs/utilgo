package nkv_test

import (
	"bytes"
	"context"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"gfx.cafe/util/go/nkv"
	"gfx.cafe/util/go/retry"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

const TestData = `
#comment
x=3
y=8
//y=13
x=6
                // y42=56
c




`

var ctx = context.TODO()

func TestParse(t *testing.T) {
	db := nkv.NewDB(func(ctx context.Context) (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte(TestData))), nil
	}, func(m map[string]string) (map[string]string, error) {
		return m, nil
	}, retry.Interval(time.Second), 0)
	v, err := db.Value(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 3 {
		t.Fatal("expected 3 values")
	}
	x, ok := v["x"]
	if !ok {
		t.Fatal("expected key x")
	}
	if x != "6" {
		t.Fatal("expected key x to have value 6")
	}
	y, ok := v["y"]
	if !ok {
		t.Fatal("expected key y")
	}
	if y != "8" {
		t.Fatal("expected key y to have value 8")
	}
	c, ok := v["c"]
	if !ok {
		t.Fatal("expected key c")
	}
	if c != "" {
		t.Fatal("expected key c to empty value")
	}
}

func TestNoModifyKvStore(t *testing.T) {
	modify := true
	db := nkv.NewDB(func(ctx context.Context) (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte(TestData))), nil
	}, func(m map[string]string) (map[string]string, error) {
		if modify {
			m["k"] = "modified"
		}
		return m, nil
	}, retry.Interval(time.Second), -1)
	v, err := db.Value(ctx)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{
		"x": "6",
		"y": "8",
		"c": "",
		"k": "modified",
	}
	assert.EqualValues(t, expected, v)
	modify = false
	delete(expected, "k")
	v, err = db.Value(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, expected, v)
}
func TestConcurrency(t *testing.T) {
	ctx := context.Background()
	db := nkv.NewDB(func(ctx context.Context) (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte(TestData))), nil
	}, func(m map[string]string) (map[string]string, error) {
		// lets see if it panics
		m["c"] = "84573189041"
		m["c"] = ""
		return m, nil
	}, retry.Interval(time.Second), time.Second)
	wg := errgroup.Group{}
	wg.SetLimit(10000)
	for i := 0; i < 100000; i++ {
		wg.Go(func() error {
			v, err := db.Value(ctx)
			if err != nil {
				return err
			}
			if len(v) != 3 {
				t.Fatal("expected 3 values")
			}
			x, ok := v["x"]
			if !ok {
				t.Fatal("expected key x")
			}
			if x != "6" {
				t.Fatal("expected key x to have value 6")
			}
			y, ok := v["y"]
			if !ok {
				t.Fatal("expected key y")
			}
			if y != "8" {
				t.Fatal("expected key y to have value 8")
			}
			c, ok := v["c"]
			if !ok {
				t.Fatal("expected key c")
			}
			if c != "" {
				t.Fatal("expected key c to empty value")
			}
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestTTLAlwaysRefresh(t *testing.T) {
	x := 0
	db := nkv.NewDB(func(ctx context.Context) (io.ReadCloser, error) {
		// should always give an empty map
		x += 1
		return io.NopCloser(bytes.NewReader(nil)), nil
	}, func(m map[string]string) (int, error) {
		return x, nil
	}, retry.Interval(time.Second), 0)
	v, err := db.Value(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Fatal("expected 1")
	}
	v, err = db.Value(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v != 2 {
		t.Fatal("expected 2")
	}
}

func TestTTLNeverRefresh(t *testing.T) {
	x := 0
	db := nkv.NewDB(func(ctx context.Context) (io.ReadCloser, error) {
		// should always give an empty map
		x += 1
		return io.NopCloser(bytes.NewReader(nil)), nil
	}, func(m map[string]string) (int, error) {
		return x, nil
	}, retry.Interval(time.Second), -1)
	v, err := db.Value(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Fatal("expected 1")
	}
	v, err = db.Value(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Fatal("expected 1")
	}
}

func TestSingleFlightWorks(t *testing.T) {
	x := atomic.Int64{}
	db := nkv.NewDB(func(ctx context.Context) (io.ReadCloser, error) {
		x.Add(1)
		time.Sleep(3 * time.Second)
		return io.NopCloser(bytes.NewReader(nil)), nil
	}, func(m map[string]string) (int, error) {
		return int(x.Load()), nil
	}, retry.Interval(time.Second), time.Second)
	var wg errgroup.Group
	var v1 int
	wg.Go(func() error {
		var err error
		v1, err = db.Value(ctx)
		return err
	})
	v2, err := db.Value(ctx)
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	if v1 != v2 {
		t.Fatal("expected same value")
	}
}
