package world

import (
	"bufio"
	"io"
	"sync"
	"time"
)

type (
	PartialPacketWriter = func(io.Writer) error

	AutoFlushingConcurrentBufferedWriter struct {
		writer   *bufio.Writer
		lock     *sync.Mutex
		interval time.Duration
		err      error
	}
)

func NewAutoFlushingConcurrentBufferedWriter(writer *bufio.Writer, interval time.Duration) *AutoFlushingConcurrentBufferedWriter {
	w := new(AutoFlushingConcurrentBufferedWriter)
	w.writer = writer
	w.lock = new(sync.Mutex)
	w.interval = interval
	return w
}

func (w *AutoFlushingConcurrentBufferedWriter) Write(ppw ...PartialPacketWriter) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.err != nil {
		return w.err
	}
	for _, v := range ppw {
		if err := v(w.writer); err != nil {
			w.err = err
			return err
		}
	}
	return nil
}

func (w AutoFlushingConcurrentBufferedWriter) autoFlush() {
	for {
		time.Sleep(w.interval)
		func() {
			w.lock.Lock()
			defer w.lock.Unlock()
			if err := w.writer.Flush(); err != nil {
				w.err = err
				return
			}
		}()
	}
}
