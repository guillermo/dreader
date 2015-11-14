package dreader

import (
	"errors"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

func checkRead(t *testing.T, i io.Reader, expectation string) error {
	ibuf := make([]byte, 15)
	n, err := i.Read(ibuf)
	if expectation != "" && err != nil {
		t.Fatal(err, expectation)
	} else {
		return err
	}

	if data := string(ibuf[:n]); data != expectation {
		t.Fatalf("Expecting to read %q but get %q", expectation, data)
	}
	return nil
}

func sendAndWaitMs(t *testing.T, w io.Writer, s string, ms int) {
	n, err := w.Write([]byte(s))
	if err != nil {
		t.Fatal(err)
	}
	if n != len(s) {
		t.Fatal(n, "!=", len(s))
	}
	time.Sleep(time.Millisecond * time.Duration(ms))

}

func TestDelayedReader(t *testing.T) {
	var wg sync.WaitGroup

	r, w := io.Pipe()

	dr := New(r, time.Millisecond*10)

	// Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		checkRead(t, dr, "hello world")
		checkRead(t, dr, " yeah !")

		checkRead(t, dr, "123")

		checkRead(t, dr, "a string that i")
		checkRead(t, dr, "s to long to re")
		checkRead(t, dr, "ad in a single ")
		checkRead(t, dr, "call to Read")

		checkRead(t, dr, "OK")
		err := checkRead(t, dr, "")
		if err != testError {
			t.Fatal(err)
		}
		err = checkRead(t, dr, "")
		if err != testError {
			t.Fatal(err)
		}
	}()

	// Producer
	wg.Add(1)
	go func() {
		defer wg.Done()

		sendAndWaitMs(t, w, "hello", 0)
		sendAndWaitMs(t, w, " world", 25)

		sendAndWaitMs(t, w, " yeah !", 25)

		sendAndWaitMs(t, w, "1", 0)
		sendAndWaitMs(t, w, "2", 1)
		sendAndWaitMs(t, w, "3", 25)

		sendAndWaitMs(t, w, "a string that is to long to read in a single call to Read", 25)

		sendAndWaitMs(t, w, "OK", 25)

		w.CloseWithError(testError)

	}()

	wg.Wait()
}

var testError = errors.New("testError")

func aTestReal(t *testing.T) {

	r := New(os.Stdin, time.Second)
	buf := make([]byte, 5)
	for {
		_, err := r.Read(buf)
		if err != nil {
			t.Fatal(err)
			break
		}
	}
}
