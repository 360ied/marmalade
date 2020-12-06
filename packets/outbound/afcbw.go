package outbound

import (
	"bufio"
	"io"
	"sync"
	"time"
)

type (
	action = func(*bufio.Writer) error

	// Auto Flushing Concurrent Buffered Writer
	AFCBW struct {
		writer   *bufio.Writer
		lock     *sync.Mutex
		interval time.Duration
		err      error
	}
)

func NewAFCBW(writer io.Writer, interval time.Duration) *AFCBW {
	w := new(AFCBW)
	w.writer = bufio.NewWriter(writer)
	w.lock = new(sync.Mutex)
	w.interval = interval
	go w.autoFlush()
	return w
}

func (w *AFCBW) do(actions ...action) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.err != nil {
		return w.err
	}
	for _, v := range actions {
		if err := v(w.writer); err != nil {
			w.err = err
			return err
		}
	}
	return nil
}

func (w AFCBW) autoFlush() {
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
