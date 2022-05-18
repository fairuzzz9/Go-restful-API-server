package utilities

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"
)

func aSlowClosingProcess(t *testing.T, ctx context.Context) {

	t.Log("A slow closing process")

	select {
	case <-ctx.Done():
		return
	default:
		time.Sleep(1 * time.Second)
	}

}

func anotherSlowClosingProcess(t *testing.T, ctx context.Context) {

	t.Log("Test another closing process with 2 seconds delay.")
	select {
	case <-ctx.Done():
		return
	default:
		time.Sleep(2 * time.Second)
	}

}

func TestGracefulAndStopWaitWrapper(t *testing.T) {

	// simulate a HTTP server
	server := &http.Server{}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Log(err)
		}

	}()

	// this is to simulate a Ctrl-C.
	// DO NOT TEST WITHOUT THIS!
	go func() {
		quit <- os.Interrupt
	}()

	Graceful(StopWaitWrapper(func(ctx context.Context) {

		aSlowClosingProcess(t, ctx)
		anotherSlowClosingProcess(t, ctx)

		t.Log("server shutdown.")

	}, 10*time.Second)) //timeout in 10 seconds.

}
