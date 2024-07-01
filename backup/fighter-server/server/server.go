package server

import (
	lib "github.com/JoTaeYang/study-go/library/server"
)

type FighterServer struct {
	*lib.SocketServer
}

func (s *FighterServer) Start() {
	s.SocketServer.Start()
}
