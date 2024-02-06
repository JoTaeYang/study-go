package server

import (
	"fmt"
	"log"
	"net"
	"syscall"

	"golang.org/x/sys/windows"
)

type SocketServer struct {
	listen windows.Handle
	hIocp  windows.Handle
}

func (s *SocketServer) Start() error {
	var err error
	var empty uintptr

	wsa := windows.WSAData{}

	err = windows.WSAStartup(uint32(0x202), &wsa)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	s.hIocp, err = windows.CreateIoCompletionPort(windows.InvalidHandle, windows.Handle(empty), 0, 0)
	if err != nil {
		return err
	}

	s.listen, err = windows.Socket(windows.AF_INET, windows.SOCK_STREAM, windows.IPPROTO_IP)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	result := &windows.AddrinfoW{}
	hints := windows.AddrinfoW{
		Family:   syscall.AF_INET,
		Socktype: windows.SOCK_STREAM,
		Protocol: syscall.IPPROTO_IP,
	}

	err = windows.GetAddrInfoW(syscall.StringToUTF16Ptr("0.0.0.0"), nil, &hints, &result)
	if nil != err {
		log.Println(err.Error())
		return err
	}

	// ip := net.IPv4zero
	// addr := [4]byte{}
	// for i := 0; i < net.IPv4len; i++ {
	// 	addr[i] = ip[i]
	// }

	err = windows.Bind(s.listen, &windows.SockaddrInet4{
		Port: 20000,
		Addr: [4]byte{0, 0, 0, 0},
	})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = windows.Listen(s.listen, windows.SOMAXCONN)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for {
		ln, err := net.Listen("tcp", ":20000")
		if nil != err {
			log.Println(err.Error())
			continue
		}
		conn, err := ln.Accept()

		handle, _ := getSocketHandle(conn)

		_, err = windows.CreateIoCompletionPort(windows.Handle(handle), s.hIocp, 0, 0)
		if nil != err {
			log.Println(err.Error())
			continue
		}

		log.Println("Hello", handle)
	}

	return nil
}

func (s *SocketServer) Accept() {
	for {
		ln, err := net.Listen("tcp", ":20000")
		if nil != err {
			log.Println(err.Error())
			continue
		}
		conn, err := ln.Accept()

		handle, _ := getSocketHandle(conn)

		_, err = windows.CreateIoCompletionPort(windows.Handle(handle), s.hIocp, 0, 0)
		if nil != err {
			log.Println(err.Error())
			continue
		}

		log.Println("Hello", handle)
	}
}

func getSocketHandle(conn net.Conn) (syscall.Handle, error) {
	// TCPConn 소켓 가져오기
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return syscall.InvalidHandle, fmt.Errorf("not a TCPConn")
	}

	// 소켓 핸들 가져오기
	file, err := tcpConn.File()
	if err != nil {
		return syscall.InvalidHandle, err
	}
	defer file.Close()

	// 소켓 핸들 얻기
	fd := file.Fd()
	connHandle := syscall.Handle(fd)

	return connHandle, nil
}
