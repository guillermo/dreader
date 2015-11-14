// Package dreader implements a delayed reader. A Reader that waits a fixed amount of time before returning.
// A DelayedReader will block Reads a specific time for more data and returns once it times out.
package dreader

import (
	"io"
	"time"
)

var DefaultSize = 1024

type DelayedReader struct {
	r    io.Reader
	cr   chan ([]byte)
	wait time.Duration
	c    chan ([]byte)
	err  error
	size int
	buf  []byte
}

// New returns a DelayedReader
func New(r io.Reader, delay time.Duration) *DelayedReader {
	dr := &DelayedReader{
		// reader
		r:  r,
		cr: make(chan ([]byte)),

		wait: delay,
		c:    make(chan ([]byte)),
		size: DefaultSize,
		buf:  nil,
	}

	go dr.readLoop()
	go dr.loop()

	return dr
}

func (dr *DelayedReader) readLoop() {
	buf := make([]byte, DefaultSize)
	for {
		n, err := dr.r.Read(buf)
		if err != nil {
			dr.err = err
			break
		}
		data := make([]byte, n, n)
		copy(data, buf[:n])
		dr.cr <- data
	}
	close(dr.cr)
}

func (dr *DelayedReader) loop() {
	var buf []byte
	var timer *time.Timer
	_ = "breakpoint"
For:

	for {
		if timer == nil {
			read, ok := <-dr.cr
			if !ok {
				break For
			}
			timer = time.NewTimer(dr.wait)
			buf = append(buf, read...)
		} else {
			select {
			case read, ok := <-dr.cr:
				if !ok {
					break For
				}
				buf = append(buf, read...)
			case <-timer.C:
				dr.c <- buf
				buf = []byte{}
				timer = nil
			}
		}

	}
	close(dr.c)
}

// Read returns any data in the internal buffer if there is any
// Then it returns any error in case there was any
// And then returns any data read from the from the underlaying reader.
// If an error is returned no more Read will be able to perform and a new DelayedReader should be created
func (dr *DelayedReader) Read(buf []byte) (n int, e error) {
	if dr.buf != nil {
		n := copy(buf, dr.buf)

		// Did we copy all the buffer?
		if len(dr.buf) == n {
			dr.buf = nil
			return n, nil
		}
		dr.buf = dr.buf[n:]

		return n, nil
	}

	if dr.err != nil {
		return 0, dr.err
	}
	data, ok := <-dr.c
	if !ok {
		return 0, dr.err
	}
	// If the data that we receive doesn't fit
	// Save it in the buffer
	n = copy(buf, data)
	if n < len(data) {
		dr.buf = data[n:]
		return n, nil
	}
	return n, nil
}
