package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JoTaeYang/study-go/fighter-server/server"
	lib "github.com/JoTaeYang/study-go/library/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	ser := server.FighterServer{
		&lib.SocketServer{},
	}

	ser.Start()

	sigs := make(chan os.Signal, 1)
	//done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	sig := <-sigs
	// 	log.Println(sig)
	// 	done <- true
	// }()

	// <-done
}
