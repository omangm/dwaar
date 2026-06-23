package tunnel

import (
	"context"
	"fmt"
	"io"
	"net"
)

func forwardLocal(ctx context.Context, rule ForwardRule, log io.Writer, counter *ByteCounter) (net.Listener, error) {
	localAddr := fmt.Sprintf("127.0.0.1:%d", rule.LocalPort)
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			fmt.Fprintf(log, "Accepted local connection from %s\n", conn.RemoteAddr())
			go proxyLocal(ctx, conn, rule, log, counter)
		}
	}()

	return listener, nil
}

func proxyLocal(ctx context.Context, local net.Conn, rule ForwardRule, log io.Writer, counter *ByteCounter) {
	defer local.Close()

	remoteAddr := fmt.Sprintf("%s:%d", rule.RemoteHost, rule.RemotePort)
	remote, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Fprintf(log, "Failed to connect to remote %s: %v\n", remoteAddr, err)
		return
	}
	defer remote.Close()
	fmt.Fprintf(log, "Connected to remote %s\n", remoteAddr)

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
		return
	case err := <-errc:
		if err != nil {
			fmt.Fprintf(log, "Connection error: %v\n", err)
		} else {
			fmt.Fprintf(log, "Connection closed\n")
		}
		return
	}
}
