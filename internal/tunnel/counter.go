package tunnel

import (
	"io"
	"sync/atomic"
)

type ByteCounter struct {
	Sent uint64
	Recv uint64
}

type CounterWriter struct {
	W       io.Writer
	Counter *uint64
}

func (cw *CounterWriter) Write(p []byte) (int, error) {
	n, err := cw.W.Write(p)
	if n > 0 && cw.Counter != nil {
		atomic.AddUint64(cw.Counter, uint64(n))
	}
	return n, err
}
