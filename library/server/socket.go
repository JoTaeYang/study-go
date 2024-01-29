package server

import (
	"log"

	"golang.org/x/sys/windows"
)

type SocketServer struct {
	listen windows.Handle
	hIocp  windows.Handle
}

func (s *SocketServer) Start() error {
	var err error
	var empty uintptr

	s.hIocp, err = windows.CreateIoCompletionPort(windows.InvalidHandle, windows.Handle(empty), 0, 0)
	if err != nil {
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
