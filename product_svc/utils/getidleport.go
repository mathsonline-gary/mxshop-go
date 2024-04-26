package utils

import "net"

func GetIdlePort() (port int, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func(lis *net.TCPListener) {
		_ = lis.Close()
	}(lis)

	return lis.Addr().(*net.TCPAddr).Port, nil
}
