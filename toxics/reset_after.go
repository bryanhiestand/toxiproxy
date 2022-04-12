package toxics

import (
	"time"
)

/*
The ResetAfterToxic closes the connection abruptly after a timeout (in ms).
The behavior of Close is set to discard any unsent/unacknowledged data by setting SetLinger to 0,
~= sets TCP RST flag and resets the connection.
If the timeout is set to 0, then the connection will be reset immediately.

*/

// The ResetAfterToxic passes data through until the timeout is reached, then closes the connection.
type ResetAfterToxic struct {
	// Time in milliseconds
	Timeout int64 `json:"timeout"`
}

func (t *ResetAfterToxic) GetBufferSize() int {
	return 1024
}

func (t *ResetAfterToxic) timeoutDelay() time.Duration {
	return time.Duration(t.Timeout) * time.Millisecond
}

func (t *ResetAfterToxic) Pipe(stub *ToxicStub) {
	for {
		select {
		case <-stub.Interrupt:
			return
		case c := <-stub.Input:
			if c == nil {
				stub.Close()
				return
			}
			stub.Output <- c
		case <-time.After(t.timeoutDelay()):
			stub.Close()
			return
		}
	}
}

func init() {
	Register("reset_after", new(ResetAfterToxic))
}
