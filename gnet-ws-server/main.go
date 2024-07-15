package main

import (
	"log"

	"github.com/JoTaeYang/study-go/gnet-ws-server/server"
	"github.com/JoTaeYang/study-go/pkg/yws"
	"github.com/panjf2000/gnet"
)

func main() {
	ws := &server.GnetWsServer{
		&yws.WsServer{},
	}

	ws.InitServer(1000)
	ws.SetMsgProc(server.MsgProc)

	log.Fatal(gnet.Serve(ws, "tcp://:30000", gnet.WithMulticore(true)))
}
