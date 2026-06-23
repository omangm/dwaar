package tunnel

import "time"

type Protocol string

const (
	ProtoTCP Protocol = "tcp"
	ProtoUDP Protocol = "udp"
)

type ForwardRule struct {
	ID         string   `yaml:"id"`
	Name       string   `yaml:"name"`
	LocalPort  int      `yaml:"local_port"`
	RemoteHost string   `yaml:"remote_host"`
	RemotePort int      `yaml:"remote_port"`
	Protocol   Protocol `yaml:"protocol"`
	Via        *SSHHost `yaml:"via,omitempty"` // nil = local forward
	AutoStart  bool     `yaml:"auto_start"`
	Tags       []string `yaml:"tags,omitempty"`
}

type SSHHost struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	IdentityFile string `yaml:"identity_file,omitempty"`
}

type TunnelStatus struct {
	Rule      ForwardRule
	State     State // Running / Stopped / Error / Connecting
	Pid       int
	BytesSent uint64
	BytesRecv uint64
	StartedAt time.Time
	Error     error
}

type State int

const (
	StateStopped State = iota
	StateConnecting
	StateRunning
	StateError
)
