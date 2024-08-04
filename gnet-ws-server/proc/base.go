package proc

import (
	"github.com/JoTaeYang/study-go/gnet-ws-server/server"
	"github.com/JoTaeYang/study-go/packet"
	"github.com/JoTaeYang/study-go/packet/stgo"
	"github.com/gobwas/ws"
)

const HEADER_LENGTH = 16

type baseAPI struct {
	Uid string
}

func InitMsgProc(server *server.GnetWsServer) {
	server.SetMsgProc(func(idx int32, msg []byte) (fr *ws.Frame, err error) {
		header := &stgo.PacketHeader{}

		packet.ByteToHeader(&msg, header)

		payload := msg[HEADER_LENGTH:]
		switch header.Pid {
		case int32(stgo.PacketID_CS_CONNECT):
		case int32(stgo.PacketID_CS_CHAT_USER_MSG):
		case int32(stgo.PacketID_CS_PONG):
		case int32(stgo.PacketID_CS_ROOM_CREATE):

		case int32(stgo.PacketID_CS_ROOM_ENTER):
			//roomproc.Enter(idx, &payload)
		case int32(stgo.PacketID_CS_ECHO):
			fr, err = Echo(idx, &payload)
		}

		return
	})
}
