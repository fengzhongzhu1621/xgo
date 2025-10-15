package server_transport

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"

	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/fengzhongzhu1621/xgo/network/ssl"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
)

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ------------------------ tcp/udp connection structures ----------------------------//

var errNotFound = errors.New("listener not found")

// ListenFd is the listener fd.
type ListenFd struct {
	Fd      uintptr
	Name    string
	Network string
	Address string
}

// GetPassedListener gets the inherited listener from parent process by network and address.
func GetPassedListener(network, address string) (interface{}, error) {
	return getPassedListener(network, address)
}

func getPassedListener(network, address string) (interface{}, error) {
	once.Do(inheritListeners)

	key := network + ":" + address
	v, ok := inheritedListenersMap.Load(key)
	if !ok {
		return nil, errNotFound
	}

	listeners := v.([]interface{})
	if len(listeners) == 0 {
		return nil, errNotFound
	}

	ln := listeners[0]
	listeners = listeners[1:]
	if len(listeners) == 0 {
		inheritedListenersMap.Delete(key)
	} else {
		inheritedListenersMap.Store(key, listeners)
	}

	return ln, nil
}

// inheritListeners stores the listener according to start listenfd and number of listenfd passed
// by environment variables.
func inheritListeners() {
	firstListenFd, err := strconv.ParseUint(os.Getenv(EnvGraceFirstFd), 10, 32)
	if err != nil {
		logging.Errorf("invalid %s, error: %v", EnvGraceFirstFd, err)
	}

	num, err := strconv.ParseUint(os.Getenv(EnvGraceRestartFdNum), 10, 32)
	if err != nil {
		logging.Errorf("invalid %s, error: %v", EnvGraceRestartFdNum, err)
	}

	for fd := firstListenFd; fd < firstListenFd+num; fd++ {
		file := os.NewFile(uintptr(fd), "")
		listener, addr, err := fileListener(file)
		file.Close()
		if err != nil {
			logging.Errorf("get file listener error: %v", err)
			continue
		}

		key := addr.Network() + ":" + addr.String()
		v, ok := inheritedListenersMap.LoadOrStore(key, []interface{}{listener})
		if ok {
			listeners := v.([]interface{})
			listeners = append(listeners, listener)
			inheritedListenersMap.Store(key, listeners)
		}
	}
}

func fileListener(file *os.File) (interface{}, net.Addr, error) {
	// Check file status.
	fin, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	// Is this a socket fd.
	if fin.Mode()&os.ModeSocket == 0 {
		return nil, nil, errFileIsNotSocket
	}

	// tcp, tcp4 or tcp6.
	if listener, err := net.FileListener(file); err == nil {
		return listener, listener.Addr(), nil
	}

	// udp, udp4 or udp6.
	if packetConn, err := net.FilePacketConn(file); err == nil {
		return packetConn, packetConn.LocalAddr(), nil
	}

	return nil, nil, errUnSupportedNetworkType
}

func getListenerFd(ln net.Listener) (*ListenFd, error) {
	sc, ok := ln.(syscall.Conn)
	if !ok {
		return nil, fmt.Errorf("getListenerFd err: %w", errUnSupportedListenerType)
	}
	fd, err := getRawFd(sc)
	if err != nil {
		return nil, fmt.Errorf("getListenerFd getRawFd err: %w", err)
	}
	return &ListenFd{
		Fd:      fd,
		Name:    "a tcp listener fd",
		Network: ln.Addr().Network(),
		Address: ln.Addr().String(),
	}, nil
}

// getRawFd acts like:
//
//	func (ln *net.TCPListener) (uintptr, error) {
//		f, err := ln.File()
//		if err != nil {
//			return 0, err
//		}
//		fd, err := f.Fd()
//		if err != nil {
//			return 0, err
//		}
//	}
//
// But it differs in an important way:
//
//	The method (*os.File).Fd() will set the original file descriptor to blocking mode as a side effect of fcntl(),
//	which will lead to indefinite hangs of Close/Read/Write, etc.
//
// References:
//   - https://github.com/golang/go/issues/29277
//   - https://github.com/golang/go/issues/29277#issuecomment-447526159
//   - https://github.com/golang/go/issues/29277#issuecomment-448117332
//   - https://github.com/golang/go/issues/43894
func getRawFd(sc syscall.Conn) (uintptr, error) {
	c, err := sc.SyscallConn()
	if err != nil {
		return 0, fmt.Errorf("sc.SyscallConn err: %w", err)
	}
	var lnFd uintptr
	if err := c.Control(func(fd uintptr) {
		lnFd = fd
	}); err != nil {
		return 0, fmt.Errorf("c.Control err: %w", err)
	}
	return lnFd, nil
}

// GetListenersFds gets listener fds.
func GetListenersFds() []*ListenFd {
	listenersFds := []*ListenFd{}
	listenersMap.Range(func(key, _ interface{}) bool {
		var (
			fd  *ListenFd
			err error
		)

		switch k := key.(type) {
		case net.Listener:
			fd, err = getListenerFd(k)
		case net.PacketConn:
			fd, err = getPacketConnFd(k)
		default:
			logging.Errorf("listener type passing not supported, type: %T", key)
			err = fmt.Errorf("not supported listener type: %T", key)
		}
		if err != nil {
			logging.Errorf("cannot get the listener fd, err: %v", err)
			return true
		}
		listenersFds = append(listenersFds, fd)
		return true
	})
	return listenersFds
}

// SaveListener saves the listener.
func SaveListener(listener interface{}) error {
	switch listener.(type) {
	case net.Listener, net.PacketConn:
		listenersMap.Store(listener, struct{}{})
	default:
		return fmt.Errorf("not supported listener type: %T", listener)
	}
	return nil
}

func mayLiftToTLSListener(ln net.Listener, opts *options.ListenServeOptions) (net.Listener, error) {
	if !(len(opts.TLSCertFile) > 0 && len(opts.TLSKeyFile) > 0) {
		return ln, nil
	}
	// Enable TLS.
	tlsConf, err := ssl.GetServerConfig(opts.CACertFile, opts.TLSCertFile, opts.TLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("tls get server config err: %w", err)
	}
	return tls.NewListener(ln, tlsConf), nil
}

func getPacketConnFd(c net.PacketConn) (*ListenFd, error) {
	sc, ok := c.(syscall.Conn)
	if !ok {
		return nil, fmt.Errorf("getPacketConnFd err: %w", errUnSupportedListenerType)
	}
	lnFd, err := getRawFd(sc)
	if err != nil {
		return nil, fmt.Errorf("getPacketConnFd getRawFd err: %w", err)
	}
	return &ListenFd{
		Fd:      lnFd,
		Name:    "a udp listener fd",
		Network: c.LocalAddr().Network(),
		Address: c.LocalAddr().String(),
	}, nil
}
