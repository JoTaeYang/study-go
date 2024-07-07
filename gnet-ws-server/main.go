package main

import (
	"log"

	"github.com/JoTaeYang/study-go/gnet-ws-server/server"
	"github.com/panjf2000/gnet"
)

func main() {
	ws := &server.WsServer{}
	log.Fatal(gnet.Serve(ws, "tcp://:30000", gnet.WithMulticore(true)))
}
