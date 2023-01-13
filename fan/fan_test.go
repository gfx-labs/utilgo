package fan_test

import (
	"testing"

	"gfx.cafe/util/go/fan"
)

func TestFanInSingleOut(t *testing.T) {

	ch := make(chan int)
	sub := fan.NewSub(ch)
	defer sub.Close()
	cnt := 128
	for i := 0; i < cnt; i++ {
		sub.Send(4)
	}
	for i := 0; i < cnt; i++ {
		<-ch
	}
}
