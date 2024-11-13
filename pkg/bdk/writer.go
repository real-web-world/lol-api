package bdk

import (
	"bufio"
	"context"
	"time"

	"go.uber.org/zap/zapcore"
)

var (
	_ zapcore.WriteSyncer = (*ConcurrentWriter)(nil)
)

type ConcurrentWriter struct {
	w      *bufio.Writer
	jobs   chan []byte
	done   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

func NewConcurrentWriter(w *bufio.Writer) *ConcurrentWriter {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConcurrentWriter{
		w:      w,
		jobs:   make(chan []byte, 1000),
		done:   make(chan struct{}),
		ctx:    ctx,
		cancel: cancel,
	}
}
func (w ConcurrentWriter) Write(msg []byte) (n int, err error) {
	if IsCtxDone(w.ctx) {
		return 0, nil
	}
	w.jobs <- msg
	return len(msg), nil
}
func (w ConcurrentWriter) Consume() {
	go func() {
		for {
			time.Sleep(time.Second)
			_ = w.w.Flush()
		}
	}()
	go func() {
		for msg := range w.jobs {
			_, _ = w.w.Write(msg)
		}
		<-w.done
	}()
}
func (w ConcurrentWriter) Flush() error {
	w.cancel()
	close(w.jobs)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-w.done:
			break loop
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
	return w.w.Flush()
}
func (w ConcurrentWriter) Sync() error {
	return w.w.Flush()
}
