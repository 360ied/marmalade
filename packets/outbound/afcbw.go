package outbound

import (
	"bufio"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"marmalade/helpers"
)

// Auto Flushing Concurrent Buffered Writer
type AFCBW struct {
	writer   *bufio.Writer
	lock     *sync.Mutex
	interval time.Duration
	err      error
}

func NewAFCBW(writer io.Writer, interval time.Duration) *AFCBW {
	w := new(AFCBW)
	w.writer = bufio.NewWriter(writer)
	w.lock = new(sync.Mutex)
	w.interval = interval
	go w.autoFlush()
	return w
}

func (w *AFCBW) Close() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.err = errors.New("AFCBW: closed")
}

func (w *AFCBW) do(actions ...helpers.Action) error {
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

func (w *AFCBW) autoFlush() {
	ok := true
	for ok {
		time.Sleep(w.interval)
		func() {
			w.lock.Lock()
			defer w.lock.Unlock()
			if w.err != nil {
				ok = false
				log.Printf("INFO/ERROR: AFCBW error: %v", w.err)
			} else if err := w.writer.Flush(); err != nil {
				w.err = err
				ok = false
				log.Printf("INFO/ERROR: AFCBW flush error: %v", err)
			}
		}()
	}
}
