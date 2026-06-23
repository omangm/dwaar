package tunnel

import (
	"bytes"
	"sync"
)

type LogBuffer struct {
	mu    sync.RWMutex
	lines []string
	max   int
}

func NewLogBuffer(max int) *LogBuffer {
	return &LogBuffer{max: max}
}

func (b *LogBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	lines := bytes.Split(p, []byte("\n"))
	for _, l := range lines {
		if len(l) > 0 {
			b.lines = append(b.lines, string(l))
		}
	}

	if len(b.lines) > b.max {
		b.lines = b.lines[len(b.lines)-b.max:]
	}
	return len(p), nil
}

func (b *LogBuffer) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var res string
	for i, l := range b.lines {
		res += l
		if i < len(b.lines)-1 {
			res += "\n"
		}
	}
	return res
}
