# dreader
Go package that implements Delayed Reader. A Reader that waits a fixed amount of time before returning.

```
PACKAGE DOCUMENTATION

package dreader
    import "."

    Package dreader implements a delayed reader. A Reader that waits a fixed
    amount of time before returning. A DelayedReader will block Reads a
    specific time for more data and returns once it times out.

    This package could be use for reading from a terminal where a escape
    sequence could be from one to several bytes with the only difference of
    the time

	r,w := io.Pipe()

	dr := dreader.New(r, time.Millisecond * 10)

	go func(){
	  buf := make([]byte, 1024)
	  n, err := dr.Read(buf)
	  if err != nil {
	     return err
	  }
	  buf[:n] == "hello world"
	}

	w.Write("hello")
	w.Write(" world")

VARIABLES

var DefaultSize = 1024

TYPES

type DelayedReader struct {
    // contains filtered or unexported fields
}

func New(r io.Reader, delay time.Duration) *DelayedReader
    New returns a DelayedReader

func (dr *DelayedReader) Read(buf []byte) (n int, e error)
    Read returns any data in the internal buffer if there is any Then it
    returns any error in case there was any And then returns any data read
    from the from the underlaying reader. If an error is returned no more
    Read will be able to perform and a new DelayedReader should be created


```
