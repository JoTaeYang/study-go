package server

import (
	"github.com/JoTaeYang/study-go/gnet-ws-server/proc"
	"github.com/JoTaeYang/study-go/packet"
	"github.com/JoTaeYang/study-go/packet/stgo"
	"github.com/JoTaeYang/study-go/pkg/yws"
)

type GnetWsServer struct {
	*yws.WsServer
}

const HEADER_LENGTH = 16

func InitMsgProc(server *GnetWsServer) {
	server.SetMsgProc(func(session *yws.WebSocketConn, msg []byte) error {
		header := &stgo.PacketHeader{}
		packet.ByteToHeader(&msg, header)

		payload := msg[HEADER_LENGTH:]
		switch header.Pid {
		case int32(stgo.PacketID_CS_CONNECT):
		case int32(stgo.PacketID_CS_CHAT_USER_MSG):
		case int32(stgo.PacketID_CS_PONG):
			proc.Pong(session, &payload)
		}
		return nil
	})
}
