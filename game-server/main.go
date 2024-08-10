package main

import (
	"log"

	"github.com/JoTaeYang/study-go/gnet-ws-server/proc"
	"github.com/JoTaeYang/study-go/gnet-ws-server/server"
	"github.com/JoTaeYang/study-go/pkg/yws"
	"github.com/panjf2000/gnet"
)

func main() {
	var err error
	err = InitConfig()
	if err != nil {
		return
	}

	ws := &server.GnetWsServer{
		&yws.WsServer{},
	}

	ws.InitServer(1000)
	go ws.GameGo()
	go ws.SendGo()
	proc.InitMsgProc(ws)

	log.Fatal(gnet.Serve(ws, "tcp://:30000", gnet.WithMulticore(true), gnet.WithNumEventLoop(4)))
}
