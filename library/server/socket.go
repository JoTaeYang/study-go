package server

import (
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

	//ip := []byte{'0', '0', '0', '0'}

	//serviceName := []byte{'0'}
	//ptr := (*uint16)(unsafe.Pointer(&ip))
	//servicePtr := (*uint16)(unsafe.Pointer(&serviceName))
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

	ip := net.IPv4zero
	addr := [4]byte{}
	for i := 0; i < net.IPv4len; i++ {
		addr[i] = ip[i]
	}
	err = windows.Bind(s.listen, &windows.SockaddrInet4{
		Port: 20000,
		Addr: addr,
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

	go s.Accept()

	return nil
}

func (s *SocketServer) Accept() {
	for {
		sock, _, err := windows.Accept(s.listen)

		if nil != err {
			log.Println(err)
			continue
		}

		log.Println("Hello", sock)
	}
}
