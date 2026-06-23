package tunnel

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func forwardSSH(ctx context.Context, rule ForwardRule, log io.Writer, counter *ByteCounter) (net.Listener, error) {
	if rule.Via == nil {
		return nil, fmt.Errorf("ssh host not provided")
	}

	auth, err := buildAuthMethods(rule.Via)
	if err != nil {
		return nil, err
	}

	cfg := &ssh.ClientConfig{
		User:            rule.Via.User,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", rule.Via.Host, rule.Via.Port)
	fmt.Fprintf(log, "Dialing SSH server %s...\n", addr)
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(log, "SSH connected.\n")

	localAddr := fmt.Sprintf("127.0.0.1:%d", rule.LocalPort)
	local, err := net.Listen("tcp", localAddr)
	if err != nil {
		client.Close()
		return nil, err
	}

	go func() {
		defer client.Close()
		for {
			conn, err := local.Accept()
			if err != nil {
				return
			}
			fmt.Fprintf(log, "Accepted local connection from %s\n", conn.RemoteAddr())
			go proxySSH(ctx, conn, client, rule, log, counter)
		}
	}()

	return local, nil
}

func proxySSH(ctx context.Context, local net.Conn, sc *ssh.Client, rule ForwardRule, log io.Writer, counter *ByteCounter) {
	remoteAddr := fmt.Sprintf("%s:%d", rule.RemoteHost, rule.RemotePort)
	remote, err := sc.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Fprintf(log, "Failed to connect to remote %s via SSH: %v\n", remoteAddr, err)
		local.Close()
		return
	}
	fmt.Fprintf(log, "Connected to remote %s via SSH\n", remoteAddr)

	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(&CounterWriter{W: remote, Counter: &counter.Sent}, local)
		errc <- err
	}()
	go func() {
		_, err := io.Copy(&CounterWriter{W: local, Counter: &counter.Recv}, remote)
		errc <- err
	}()

	select {
	case <-ctx.Done():
		local.Close()
		remote.Close()
		return
	case err := <-errc:
		if err != nil {
			fmt.Fprintf(log, "Connection error: %v\n", err)
		} else {
			fmt.Fprintf(log, "Connection closed\n")
		}
		local.Close()
		remote.Close()
		return
	}
}

func buildAuthMethods(host *SSHHost) ([]ssh.AuthMethod, error) {
	var methods []ssh.AuthMethod
	if host.IdentityFile != "" {
		key, err := os.ReadFile(host.IdentityFile)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}
	return methods, nil
}
