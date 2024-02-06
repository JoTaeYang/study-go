package main

import (
	"log"

	"github.com/JoTaeYang/study-go/fighter-server/server"
	lib "github.com/JoTaeYang/study-go/library/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	ser := server.FighterServer{
		&lib.SocketServer{},
	}

	ser.Start()
}
