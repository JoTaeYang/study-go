package main

import (
	"log"

	"github.com/panjf2000/gnet"
)

type echoServer struct {
	*gnet.EventServer
}

func (es *echoServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("TCP Echo server started on %s\n", srv.Addr.String())
	return
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// 클라이언트로부터 받은 데이터를 그대로 반환합니다 (에코)
	out = frame
	return
}

func main() {
	echo := &echoServer{}
	// TCP 서버를 시작합니다.
	log.Fatal(gnet.Serve(echo, "tcp://:9000", gnet.WithMulticore(true)))
}
