package graceful

import (
	"context"
	"log"
	"syscall"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("sending SIGINT")
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	log.Println("press ctrl-c to stop")
	err := Handler(5*time.Second, func(context.Context, <-chan struct{}) error {
		return nil
	}, func(ctx context.Context) error {
		for i := 4; i > 0; i-- {
			log.Printf("shutting down in %d\n", i)
			time.Sleep(1 * time.Second)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error")

	}
}
