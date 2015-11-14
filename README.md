# dreader
Go package that implements Delayed Reader. A Reader that waits a fixed amount of time before returning.


```
PACKAGE DOCUMENTATION

package dreader
    import "."

    Package dreader implements a delayed reader. A Reader that waits a fixed
    amount of time before returning. A DelayedReader will block Reads a
    specific time for more data and returns once it times out.

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

