package tunnel

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type StatusEvent struct {
	ID     string
	Status TunnelStatus
}

type Manager struct {
	mu        sync.RWMutex
	tunnels   map[string]*activeTunnel
	eventsCh  chan StatusEvent
	globalLog *LogBuffer
}

type activeTunnel struct {
	rule     ForwardRule
	cancel   context.CancelFunc
	listener net.Listener
	logBuf   *LogBuffer
	counter  *ByteCounter
	started  time.Time
}

func NewManager() *Manager {
	return &Manager{
		tunnels:   make(map[string]*activeTunnel),
		eventsCh:  make(chan StatusEvent, 100),
		globalLog: NewLogBuffer(500),
	}
}

func (m *Manager) Start(rule ForwardRule) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tunnels[rule.ID]; ok {
		return fmt.Errorf("tunnel %s already running", rule.ID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	logBuf := NewLogBuffer(200)
	counter := &ByteCounter{}
	started := time.Now()

	var listener net.Listener
	var err error

	fmt.Fprintf(logBuf, "Starting tunnel %s...\n", rule.Name)
	fmt.Fprintf(m.globalLog, "[%s] Starting tunnel...\n", rule.Name)

	if rule.Via == nil {
		listener, err = forwardLocal(ctx, rule, logBuf, counter)
	} else {
		listener, err = forwardSSH(ctx, rule, logBuf, counter)
	}

	if err != nil {
		cancel()
		fmt.Fprintf(logBuf, "Failed to start tunnel: %v\n", err)
		fmt.Fprintf(m.globalLog, "[%s] Failed to start tunnel: %v\n", rule.Name, err)
		m.eventsCh <- StatusEvent{
			ID: rule.ID,
			Status: TunnelStatus{
				Rule:  rule,
				State: StateError,
				Error: err,
			},
		}
		return err
	}

	fmt.Fprintf(logBuf, "Tunnel started successfully.\n")
	fmt.Fprintf(m.globalLog, "[%s] Tunnel started successfully.\n", rule.Name)

	m.tunnels[rule.ID] = &activeTunnel{
		rule:     rule,
		cancel:   cancel,
		listener: listener,
		logBuf:   logBuf,
		counter:  counter,
		started:  started,
	}

	m.eventsCh <- StatusEvent{
		ID: rule.ID,
		Status: TunnelStatus{
			Rule:  rule,
			State: StateRunning,
		},
	}

	return nil
}

func (m *Manager) Stop(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t, ok := m.tunnels[id]; ok {
		t.cancel()
		if t.listener != nil {
			t.listener.Close()
		}
		fmt.Fprintf(t.logBuf, "Tunnel stopped.\n")
		fmt.Fprintf(m.globalLog, "[%s] Tunnel stopped.\n", t.rule.Name)
		delete(m.tunnels, id)
		
		m.eventsCh <- StatusEvent{
			ID: id,
			Status: TunnelStatus{
				Rule:  t.rule,
				State: StateStopped,
			},
		}
	}
	return nil
}

func (m *Manager) Status(id string) TunnelStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if t, ok := m.tunnels[id]; ok {
		return TunnelStatus{
			Rule:      t.rule,
			State:     StateRunning,
			BytesSent: atomic.LoadUint64(&t.counter.Sent),
			BytesRecv: atomic.LoadUint64(&t.counter.Recv),
			StartedAt: t.started,
		}
	}
	return TunnelStatus{State: StateStopped}
}

func (m *Manager) GetLogs(id string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if t, ok := m.tunnels[id]; ok {
		return t.logBuf.String()
	}
	return "No logs available."
}

func (m *Manager) GetGlobalLogs() string {
	return m.globalLog.String()
}

func (m *Manager) Events() <-chan StatusEvent {
	return m.eventsCh
}
