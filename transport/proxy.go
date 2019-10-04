package transport

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
)

// Proxy copies data between the given network connections.
func Proxy(src net.Conn, dst net.Conn) {
	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			fmt.Printf("Dqlite proxy TLS -> Unix: %v\n", err)
		}
		src.Close()
		dst.Close()
	}()

	go func() {
		_, err := io.Copy(src, dst)
		if err != nil {
			fmt.Printf("Dqlite proxy Unix -> TLS: %v\n", err)
		}
		src.Close()
		dst.Close()
	}()
}

// Socketpair returns a pair of connected unix sockets.
func Socketpair() (net.Conn, net.Conn, error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, nil, err
	}

	c1, err := fdToFileConn(fds[0])
	if err != nil {
		return nil, nil, err
	}

	c2, err := fdToFileConn(fds[1])
	if err != nil {
		c1.Close()
		return nil, nil, err
	}

	return c1, c2, err
}

func fdToFileConn(fd int) (net.Conn, error) {
	f := os.NewFile(uintptr(fd), "")
	defer f.Close()
	return net.FileConn(f)
}
