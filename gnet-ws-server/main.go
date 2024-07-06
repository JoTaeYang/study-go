package main

import (
	"log"

	"github.com/panjf2000/gnet"
)

func main() {
	ws := &server.wsServer{}
	log.Fatal(gnet.Serve(ws, "tcp://:30000", gnet.WithMulticore(true)))
}
