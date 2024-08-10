package main

import (
	"log"

	"github.com/JoTaeYang/study-go/game-server/server"
	"github.com/panjf2000/gnet"
)

func main() {
	var err error
	err = InitConfig()
	if err != nil {
		return
	}

	s := server.NewGameServer()

	s.InitServer(1000)
	//go s.GameGo()
	go s.SendGo()
	server.InitMsgProc(s)

	log.Fatal(gnet.Serve(s, "tcp://:20000", gnet.WithMulticore(true), gnet.WithNumEventLoop(4)))
}
