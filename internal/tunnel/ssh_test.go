package tunnel

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// directTCPIPPayload mirrors the RFC 4254 direct-tcpip channel payload.
type directTCPIPPayload struct {
	DestAddr string
	DestPort uint32
	OrigAddr string
	OrigPort uint32
}

// startTestSSHServer runs an in-process SSH server that accepts the given
// user/password and proxies direct-tcpip channels to the requested address.
func startTestSSHServer(t *testing.T, user, pass string) string {
	t.Helper()

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, p []byte) (*ssh.Permissions, error) {
			if c.User() == user && string(p) == pass {
				return nil, nil
			}
			return nil, fmt.Errorf("access denied")
		},
	}
	cfg.AddHostKey(signer)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go handleSSHConn(conn, cfg)
		}
	}()

	_, port, _ := net.SplitHostPort(ln.Addr().String())
	return port
}

func handleSSHConn(conn net.Conn, cfg *ssh.ServerConfig) {
	sconn, chans, reqs, err := ssh.NewServerConn(conn, cfg)
	if err != nil {
		conn.Close()
		return
	}
	defer sconn.Close()
	go ssh.DiscardRequests(reqs)

	for ch := range chans {
		if ch.ChannelType() != "direct-tcpip" {
			ch.Reject(ssh.UnknownChannelType, "only direct-tcpip supported")
			continue
		}
		var p directTCPIPPayload
		if err := ssh.Unmarshal(ch.ExtraData(), &p); err != nil {
			ch.Reject(ssh.ConnectionFailed, "bad payload")
			continue
		}
		remote, err := net.Dial("tcp", net.JoinHostPort(p.DestAddr, fmt.Sprint(p.DestPort)))
		if err != nil {
			ch.Reject(ssh.ConnectionFailed, err.Error())
			continue
		}
		chConn, reqs, err := ch.Accept()
		if err != nil {
			remote.Close()
			continue
		}
		go ssh.DiscardRequests(reqs)
		go func() {
			io.Copy(chConn, remote)
			chConn.Close()
		}()
		go func() {
			io.Copy(remote, chConn)
			remote.Close()
		}()
	}
}

// startEchoServer returns the port of a TCP server that echoes back one message.
func startEchoServer(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c)
			}(conn)
		}
	}()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	return port
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestForwardSSHPasswordAuth(t *testing.T) {
	echoPort := startEchoServer(t)
	sshPort := startTestSSHServer(t, "dwaar", "s3cret")
	localPort := freePort(t)

	rule := ForwardRule{
		ID:         "test",
		Name:       "test",
		LocalPort:  localPort,
		RemoteHost: "127.0.0.1",
		RemotePort: atoi(echoPort),
		Via: &SSHHost{
			Host:     "127.0.0.1",
			Port:     atoi(sshPort),
			User:     "dwaar",
			Password: "s3cret",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener, err := forwardSSH(ctx, rule, io.Discard, &ByteCounter{})
	if err != nil {
		t.Fatalf("forwardSSH: %v", err)
	}
	defer listener.Close()

	conn, err := net.DialTimeout("tcp", listener.Addr().String(), 2*time.Second)
	if err != nil {
		t.Fatalf("dial local port: %v", err)
	}
	defer conn.Close()

	msg := []byte("dwaar-ping")
	if _, err := conn.Write(msg); err != nil {
		t.Fatalf("write: %v", err)
	}
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, len(msg))
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(buf) != string(msg) {
		t.Fatalf("echo mismatch: got %q want %q", buf, msg)
	}
}

func TestForwardSSHBadPassword(t *testing.T) {
	echoPort := startEchoServer(t)
	sshPort := startTestSSHServer(t, "dwaar", "s3cret")

	rule := ForwardRule{
		ID:         "test",
		Name:       "test",
		LocalPort:  freePort(t),
		RemoteHost: "127.0.0.1",
		RemotePort: atoi(echoPort),
		Via: &SSHHost{
			Host:     "127.0.0.1",
			Port:     atoi(sshPort),
			User:     "dwaar",
			Password: "wrong",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if _, err := forwardSSH(ctx, rule, io.Discard, &ByteCounter{}); err == nil {
		t.Fatal("expected auth error, got nil")
	}
}

func TestBuildAuthMethodsNone(t *testing.T) {
	if _, err := buildAuthMethods(&SSHHost{Host: "h", User: "u"}); err == nil {
		t.Fatal("expected error when no auth method configured")
	}
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
